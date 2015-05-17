package plugins

import (
	"time"
)

type Plugin interface {
	Init(map[string]chan []byte)
	Execute()
	Name() string
	IsRunning() bool
	DataChannel() chan [][]byte
	Interval() time.Duration
}

type BasePlugin struct {
	Pool     map[string]chan []byte
	out      chan [][]byte
	interval time.Duration
}

func (pl *BasePlugin) Init(pool map[string]chan []byte) {
	pl.Pool = pool
	pl.out = make(chan [][]byte)
	pl.interval = time.Second
}

func (pl *BasePlugin) DataChannel() chan [][]byte {
	return pl.out
}

func (pl *BasePlugin) Interval() time.Duration {
	return pl.interval
}
