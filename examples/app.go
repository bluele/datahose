package main

import (
	"github.com/bluele/datahose"
	"github.com/bluele/datahose/examples/plugins/times"
	"log"
)

const (
	addr = ":8000"
)

func main() {
	hose := datahose.New()
	hose.Register("/times", times.New())
	log.Println("serving at " + addr)
	hose.Serve(addr)
}
