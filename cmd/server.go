package main

import (
	"github.com/takahawk/shadownet/gateway"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/storages"
)

const DefaultDatabaseFilename = "shadownet.db"

func main() {
	// TODO: set port through options
	port := 1337

	logger := logger.NewZerologLogger(logger.NewZerologLoggerConfig())
	logger.Info("Hello world")
	storage, err := storages.NewSqliteStorage(DefaultDatabaseFilename, logger)
	if err != nil {
		return
	}
	gateway := gateway.NewShadowGateway(logger, storage)
	gateway.Start(port)
}
