package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/encryptors"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type builtinResolver struct {
	downloaderDict map[string] func(params... string) (downloaders.Downloader, error)
	transformerDict map[string] func(params... string) (transformers.Transformer, error)
	encryptorDict map[string] func(params... string) (encryptors.Encryptor, error)
	uploaderDict map[string] func(params... string) (uploaders.Uploader, error)
}

func NewBuiltinResolver() Resolver {
	return &builtinResolver{
		downloaderDict: map[string] func(params... string) (downloaders.Downloader, error) {
			downloaders.PastebinDownloaderName: func(params... string) (downloaders.Downloader, error) { 
				if len(params) != 0 {
					return nil, errors.New("pastebin downloader doesn't accept any params")
				}
				return downloaders.NewPastebinDownloader(), nil
			},
		},

		transformerDict: map[string] func(params... string) (transformers.Transformer, error) {
			transformers.Base64TransformerName: func(params... string) (transformers.Transformer, error) {
				if len(params) != 0 {
					return nil, errors.New("base64 transformer doesn't accept any params")
				}

				return transformers.NewBase64Transformer(), nil
			},
		},

		encryptorDict: map[string] func(params... string) (encryptors.Encryptor, error) {
			encryptors.AESEncryptorName: func(params... string) (encryptors.Encryptor, error) {
				return encryptors.NewAESEncryptorWithParams(params...)
			},
		},

		uploaderDict: map[string] func(params... string) (uploaders.Uploader, error) {
			uploaders.PastebinUploaderName: func(params... string) (uploaders.Uploader, error) {
				return uploaders.NewPastebinUploaderWithParams(params...)
			},
		},
	}
}

func (br *builtinResolver) Resolve(urlPart string) (common.Nameable, error) {
	return nil, nil
}


func (br *builtinResolver) ResolveDownloader(name string, params... string) (downloaders.Downloader, error) {
	if downloaderFactory, ok := br.downloaderDict[name]; ok {
		return downloaderFactory(params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in downloader with name %s", name))
}

func (br *builtinResolver) ResolveTransformer(name string, params... string) (transformers.Transformer, error) {
	if transformerFactory, ok := br.transformerDict[name]; ok {
		return transformerFactory(params...)
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in transformer with name %s", name))
}

func (br *builtinResolver) ResolveEncryptor(name string, params... string) (encryptors.Encryptor, error) {
	if encryptorFactory, ok := br.encryptorDict[name]; ok {
		return encryptorFactory(params...)
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in encryptor with name %s", name))
}

func (br *builtinResolver) ResolveUploader(name string, params... string) (uploaders.Uploader, error) {
	if uploaderFactory, ok := br.uploaderDict[name]; ok {
		return uploaderFactory(params...)
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in uploader with name %s", name))
}