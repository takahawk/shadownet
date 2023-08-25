package main

import (
	"fmt"
	"github.com/takahawk/shadownet/uploaders"
)

func main() {
	devKey := "" // TODO: set
	webPage := []byte("<HTML><BODY><b>Bold text</b><br><i>italic text</i></BODY></HTML>")
	
	pastebinUploader := uploaders.NewPastebinUploader(devKey)
	id, err := pastebinUploader.Upload(webPage)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		return
	}
	fmt.Printf("ID: %s", id)
}