package main

import (
	"fmt"
	"net/http"

	"github.com/takahawk/shadownet/resolvers"
)

func gateway(w http.ResponseWriter, req *http.Request) {
	resolver := resolvers.NewBuiltinResolver()
	downloader, _ := resolver.ResolveDownloader("pastebin", []byte("yHWR5RQr"))
	content, err := downloader.Download()
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%+v", err))
	}
	
	fmt.Fprintf(w, string(content))
}


func main() {
	
	testtext := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque congue nisi orci, in convallis eros faucibus nec. Cras bibendum elit nisi, et euismod justo dignissim vitae. Donec rutrum tortor euismod ullamcorper pharetra. Nulla nec est metus. Donec ut luctus metus. Maecenas ac mauris et mauris consectetur gravida. Ut vitae laoreet arcu. Donec venenatis tortor non nunc tristique, a rhoncus justo vestibulum. "
	resolver := resolvers.NewBuiltinResolver()
	encryptor, err := resolver.ResolveEncryptor("aes", []byte("thereisnospoonthereisnospoonther"), []byte("abcdefghabcdefgh"))
	if err != nil {
		fmt.Printf("error: %+v", err)
		return
	}

	transformer, _ := resolver.ResolveTransformer("base64")
	
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