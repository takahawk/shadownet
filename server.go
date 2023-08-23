package main

import (
	"fmt"
	"net/http"

	"github.com/takahawk/shadownet/encryptors"
	"github.com/takahawk/shadownet/storages"
	"github.com/takahawk/shadownet/transformers"
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
	testtext := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque congue nisi orci, in convallis eros faucibus nec. Cras bibendum elit nisi, et euismod justo dignissim vitae. Donec rutrum tortor euismod ullamcorper pharetra. Nulla nec est metus. Donec ut luctus metus. Maecenas ac mauris et mauris consectetur gravida. Ut vitae laoreet arcu. Donec venenatis tortor non nunc tristique, a rhoncus justo vestibulum. "
	encryptor, err := encryptors.NewAESEncryptor([]byte("thereisnospoonthereisnospoonther"), []byte("abcdefghabcdefgh"))
	if err != nil {
		fmt.Printf("%+v", err)
	}
	transformer := transformers.NewBase64Transformer()
	
	cipher, err := encryptor.Encrypt([]byte(testtext))
	fmt.Printf("Source text: %s\n", testtext)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("Before base (len=%d): %s\n", len(cipher), cipher)
	cipher, err = transformer.ForwardTransform(cipher)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("Encrypted text (len=%d): %s\n",  len(cipher), string(cipher))

	decrypted, err := transformer.ReverseTransform(cipher)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("Unbase (len=%d): %s\n", len(decrypted), decrypted)
	decrypted, err = encryptor.Decrypt(decrypted)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	
	fmt.Printf("Decrypted text: %s\n", string(decrypted))


	port := 1337
	http.HandleFunc("/", gateway)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}