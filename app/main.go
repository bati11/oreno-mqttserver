package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/bati11/oreno-mqtt/mqtt"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	mqtt.Run()
}
