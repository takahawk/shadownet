package main

import (
	"fmt"
	"net/http"

	"github.com/takahawk/shadownet/storages"
)

func gateway(w http.ResponseWriter, req *http.Request) {
	downloader := storages.NewPastebinDownloader()
	content, err := downloader.Download("yHWR5RQr")
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%+v", err))
	}
	fmt.Fprintf(w, content)
}


func main() {
	port := 1337
	http.HandleFunc("/", gateway)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}