package main

import (
	"fmt"

	"os"
//	"github.com/pborman/getopt"

	"github.com/takahawk/shadownet/pipelines"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
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

	data := []byte("<html><body><b>Bold text</b><br><i>Italic text</i><br>Plain text</body></html")
	pipeline := pipelines.NewUploadPipeline()
	encryptor, _ := transformers.NewAESEncryptor([]byte("thereisnospoonthereisnospoonther"), []byte("abcdefghabcdefgh"))
	err := pipeline.AddSteps(
		encryptor, 
		transformers.NewBase64Transformer(), 
		uploaders.NewPastebinUploader(""))
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(-1)
	}
	url, err := pipeline.Upload(data)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(-1)
	}
	fmt.Printf("URL: %+v", url)
}