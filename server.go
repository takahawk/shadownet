package main

import (
	"fmt"
	"net/http"

	"github.com/takahawk/shadownet/encryptors"
	"github.com/takahawk/shadownet/storages"
)

func gateway(w http.ResponseWriter, req *http.Request) {
	testtext := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque congue nisi orci, in convallis eros faucibus nec. Cras bibendum elit nisi, et euismod justo dignissim vitae. Donec rutrum tortor euismod ullamcorper pharetra. Nulla nec est metus. Donec ut luctus metus. Maecenas ac mauris et mauris consectetur gravida. Ut vitae laoreet arcu. Donec venenatis tortor non nunc tristique, a rhoncus justo vestibulum. "
	encryptor := encryptors.NewAESEncryptor("thereisnospoonthereisnospoonther")
	cipher, err := encryptor.Encrypt(testtext)
	fmt.Printf("Source text: %s\n", testtext)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("Encrypted text: %s\n", cipher)

	decrypted, err := encryptor.Decrypt(cipher)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("Decrypted text: %s\n", decrypted)

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