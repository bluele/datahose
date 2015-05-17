package connections

import (
	"github.com/bluele/datahose/lib/plugins"
	"github.com/nu7hatch/gouuid"
	"log"
	"sync"
)

type ConnectionPool struct {
	// plugin => conn_id => queue
	pool  map[string]map[string]chan []byte
	mutex sync.RWMutex
}

type Connection struct {
	ID     string
	Path   string
	Queue  chan []byte
	Plugin plugins.Plugin
}

var Pool *ConnectionPool

func NewConnectionPool() *ConnectionPool {
	cp := &ConnectionPool{}
	cp.pool = make(map[string]map[string]chan []byte)
	return cp
}

func (cp *ConnectionPool) Get(pl plugins.Plugin) *Connection {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	path := pl.Name()
	conn := &Connection{
		ID:     id.String(),
		Path:   path,
		Queue:  make(chan []byte, 10),
		Plugin: pl,
	}
	cp.mutex.Lock()
	cp.pool[path][conn.ID] = conn.Queue
	cp.mutex.Unlock()
	return conn
}

func (cp *ConnectionPool) Release(conn *Connection) {
	cp.mutex.Lock()
	close(conn.Queue)
	delete(cp.pool[conn.Plugin.Name()], conn.ID)
	cp.mutex.Unlock()
}

func (cp *ConnectionPool) PutMessage(path string, msg []byte) bool {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	conns, ok := cp.pool[path]
	if !ok {
		return false
	}
	for id, conn := range conns {
		select {
		case conn <- msg:
		default:
			log.Println("Cannot push message to connection: " + id)
		}
	}
	return true
}

func (cp *ConnectionPool) SetConnectionMap(path string, pool map[string]chan []byte) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.pool[path] = pool
}

func init() {
	Pool = NewConnectionPool()
}
