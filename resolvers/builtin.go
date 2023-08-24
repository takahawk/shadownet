package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/encryptors"
)

type builtinDownloaderResolver struct {
	dict map[string] downloaders.Downloader
}
type builtinTransformerResolver struct {
	dict map[string] transformers.Transformer
}
type builtinEncryptorResolver struct {
	dict map[string] encryptors.Encryptor
}

func NewBuiltinDownloaderResolver() DownloaderResolver {
	return &builtinDownloaderResolver{
		dict: map[string] downloaders.Downloader {
			"pastebin": downloaders.NewPastebinDownloader(),
		},
	}
}

func NewBuiltinTransformerResolver() TransformerResolver {
	return &builtinTransformerResolver{
		dict: map[string] transformers.Transformer {
			"base64": transformers.NewBase64Transformer(),
		},
	}
}

func NewBuiltinEncryptorResolver() EncryptorResolver {
	aesEncryptor, _ := encryptors.NewAESEncryptor([]byte("thereisnospoonthereisnospoonther"), []byte("abcdefghabcdefgh"))
	return &builtinEncryptorResolver{
		dict: map[string] encryptors.Encryptor {
			// TODO: move key and iv to interface from constructor
			"aes": aesEncryptor,
		},
	}
}

func (bdr *builtinDownloaderResolver) ResolveDownloader(id string) (downloaders.Downloader, error) {
	if downloader, ok := bdr.dict[id]; ok {
		return downloader, nil
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in downloader with ID %s", id))
}

func (btr *builtinTransformerResolver) ResolveTransformer(id string) (transformers.Transformer, error) {
	if transformer, ok := btr.dict[id]; ok {
		return transformer, nil
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in transformer with ID %s", id))
}

func (ber *builtinEncryptorResolver) ResolveEncryptor(id string) (encryptors.Encryptor, error) {
	if encryptor, ok := ber.dict[id]; ok {
		return encryptor, nil
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in encryptor with ID %s", id))
}