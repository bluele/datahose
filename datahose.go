package datahose

import (
	"github.com/bluele/datahose/lib/connections"
	"github.com/bluele/datahose/lib/plugins"
	"github.com/bluele/datahose/lib/tasks"
	"net/http"
	"time"
)

var (
	DefaultInterval = time.Second
)

type Hose struct {
	Mux          *http.ServeMux
	TaskInterval time.Duration
}

func New() *Hose {
	return &Hose{
		Mux:          http.NewServeMux(),
		TaskInterval: DefaultInterval,
	}
}

func (ho *Hose) Register(prefix string, pl plugins.Plugin) {
	ho.Mux.HandleFunc(prefix, MakeHandler(pl))
}

func (ho *Hose) Serve(addr string) {
	tasks.Manager.Execute()
	http.ListenAndServe(addr, ho.Mux)
}

func handler(pl plugins.Plugin) func(http.ResponseWriter, *http.Request) {
	var msg []byte
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		notify := w.(http.CloseNotifier).CloseNotify()
		conn := connections.Pool.Get(pl)
		defer connections.Pool.Release(conn)
		flusher := w.(http.Flusher)
	Main:
		for {
			select {
			case msg = <-conn.Queue:
				w.Write(append(msg, 10))
				flusher.Flush()
			case <-notify:
				break Main
			}
		}
	}
}

func MakeHandler(pl plugins.Plugin) func(http.ResponseWriter, *http.Request) {
	tasks.Manager.CreateTask(pl)
	return handler(pl)
}
