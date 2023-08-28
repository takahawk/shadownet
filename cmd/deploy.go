package main

import (
	"fmt"

	"os"
//	"github.com/pborman/getopt"

	// "github.com/takahawk/shadownet/pipelines"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	// "github.com/takahawk/shadownet/uploaders"
	"github.com/takahawk/shadownet/url"
)

func main() {
	// optUploader := getopt.StringLong("uploader", 'u', "", "Uploader name")
	// optParameters := getopt.StringLong("params", 'p', "", "Uploader parameters")
	// optTransformer := getopt.StringLong("transformer", 't', "", "Transformer name")
	// optEncryptor := getopt.StringLong("encryptor", 'e', "", "Encryptor name")
	// optFile := getopt.StringLong("file", 'f', "", "File to upload")
	// optHelp := getopt.BoolLong("help", 'h', "Help")
	// // devKey := "" // TODO: set
	
	// getopt.Parse()

	// if optHelp {
	// 	getopt.Usage()
	// 	os.Exit(0)
	// }
	// if optName == "" {
		
	// 	fmt.Printf("'uploader' parameter is required")
	// 	os.Exit(-1)
	// }

	// if optFile == "" {
	// 	getopt.Usage()
	// 	fmt.Printf("'file' parameter is required")
	// 	os.Exit(-1)
	// }

	// data := []byte("<html><body><b>Bold text</b><br><i>Italic text</i><br>Plain text</body></html")
	// pipeline := pipelines.NewUploadPipeline()
	// encryptor, _ := transformers.NewAESEncryptor([]byte("thereisnospoonthereisnospoonther"), []byte("abcdefghabcdefgh"))
	// err := pipeline.AddSteps(
	// 	encryptor, 
	// 	transformers.NewBase64Transformer(), 
	// 	uploaders.NewPastebinUploader(""))
	// if err != nil {
	// 	fmt.Printf("Error: %+v", err)
	// 	os.Exit(-1)
	// }
	// shadowURL, err := pipeline.Upload(data)
	// if err != nil {
	// 	fmt.Printf("Error: %+v", err)
	// 	os.Exit(-1)
	// }
	// fmt.Printf("URL: %+v\n", shadowURL)
	shadowURL := "ZG93bl9wYXN0ZWJpbjpjRzU1Y1Zaa2RuST0=.dHJhbnNfYmFzZTY0Og==.dHJhbnNfYWVzOmRHaGxjbVZwYzI1dmMzQnZiMjUwYUdWeVpXbHpibTl6Y0c5dmJuUm9aWEk9LFlXSmpaR1ZtWjJoaFltTmtaV1puYUE9PQ=="

	urlHandler := url.NewUrlHandler()
	components, err := urlHandler.GetDownloadComponents(shadowURL)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(-1)
	}

	var data []byte
	for _, component := range components {
		switch component := component.(type) {
		case downloaders.Downloader:
			data, err = component.Download()
		case transformers.Transformer:
			data, err = component.ReverseTransform(data)
		}

		if err != nil {
			fmt.Printf("Error: %+v", err)
			os.Exit(-1)
		}
	}
	fmt.Printf("Data: %+v", string(data))
}