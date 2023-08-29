package main

import (
	"github.com/takahawk/shadownet/gateway"
)


func main() {
	// TODO: set port through options
	port := 1337
	gateway := gateway.NewShadowGateway()
	gateway.Start(port)
}