package main

import (
	"github.com/takahawk/shadownet/gateway"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/storages"
	"github.com/takahawk/shadownet/uploaders"
)

const DefaultDatabaseFilename = "shadownet.db"
const ShadowNetPort = 10176

func main() {
	// TODO: set port through options

	logger := logger.NewZerologLogger(logger.NewZerologLoggerConfig())
	logger.Info("Hello world")
	storage, err := storages.NewSqliteStorage(DefaultDatabaseFilename, logger)
	if err != nil {
		return
	}

	data := []byte("<html><body><b>Bold text</b><br><i>Italic text</i><br>Plain text</body></html")
	uploader := uploaders.NewDropboxUploader(logger, "sl.BmDvCi4cSKqTNJTN-QSdi-SgjAMn_oIP3pa5D_NhjcMkcdj69MPAYs7oAP6jbmPwklQIyhYce-vUpZZ5IrD3GBzScR6yXwCvdJ_SwDRbxJ4haBHaenf6DDn2a2hDH1N4AeZBEpeoRuekc0M")

	uploader.Upload(data)
	gateway := gateway.NewShadowGateway(logger, storage)
	gateway.Start(ShadowNetPort)
}
