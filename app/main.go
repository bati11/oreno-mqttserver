package main

import (
	"flag"
	_ "net/http/pprof"

	"github.com/bati11/oreno-mqtt/mqtt"
)

func main() {
	ws := flag.Bool("w", false, "open websocket port (default \"false\")")
	flag.Parse()
	mqtt.Run(*ws)
}
