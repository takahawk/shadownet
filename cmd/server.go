package main

import (
	"github.com/takahawk/shadownet/gateway"
	"github.com/takahawk/shadownet/logger"
)

func main() {
	// TODO: set port through options
	port := 1337

	logger := logger.NewZerologLogger(logger.NewZerologLoggerConfig())
	logger.Info("Hello world")
	gateway := gateway.NewShadowGateway(logger)
	gateway.Start(port)
}
