package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/encryptors"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
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

type builtinUploaderResolver struct {
	dict map[string] func(params interface{}) (uploaders.Uploader, error)
}

type PastebinUploaderParams struct {
	devKey string
}

func NewBuiltinDownloaderResolver() DownloaderResolver {
	pastebinDownloader := downloaders.NewPastebinDownloader()
	return &builtinDownloaderResolver{
		dict: map[string] downloaders.Downloader {
			pastebinDownloader.Name(): pastebinDownloader,
		},
	}
}

func NewBuiltinTransformerResolver() TransformerResolver {
	base64Transformer := transformers.NewBase64Transformer()
	return &builtinTransformerResolver{
		dict: map[string] transformers.Transformer {
			base64Transformer.Name(): base64Transformer,
		},
	}
}

func NewBuiltinEncryptorResolver() EncryptorResolver {
	aesEncryptor := encryptors.NewAESEncryptor()
	return &builtinEncryptorResolver{
		dict: map[string] encryptors.Encryptor {
			// TODO: move key and iv to interface from constructor
			aesEncryptor.Name(): aesEncryptor,
		},
	}
}

func NewBuiltinUploaderResolver() UploaderResolver {
	return &builtinUploaderResolver{
		dict: map[string] func(params interface{}) (uploaders.Uploader, error) {
			// TODO: move key and iv to interface from constructor
			uploaders.NewPastebinUploader("").Name(): func(params interface{}) (uploaders.Uploader, error) {
				switch p := params.(type) {
				case PastebinUploaderParams:
					return uploaders.NewPastebinUploader(p.devKey), nil
				default:
					return nil, errors.New(fmt.Sprintf("unknown param type %T. should be PastebinUploaderParams", p))
				}
			},
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

func (bur *builtinUploaderResolver) ResolveUploader(id string) (func(params interface{}) (uploaders.Uploader, error), error) {
	if uploader, ok := bur.dict[id]; ok {
		return uploader, nil
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in uploader with ID %s", id))
}