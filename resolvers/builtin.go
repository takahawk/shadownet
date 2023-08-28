package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type builtinResolver struct {
	downloaderDict map[string] func(params... []byte) (downloaders.Downloader, error)
	transformerDict map[string] func(params... []byte) (transformers.Transformer, error)
	uploaderDict map[string] func(params... []byte) (uploaders.Uploader, error)
}

func NewBuiltinResolver() Resolver {
	return &builtinResolver{
		downloaderDict: map[string] func(params... []byte) (downloaders.Downloader, error) {
			downloaders.PastebinDownloaderName: func(params... []byte) (downloaders.Downloader, error) { 
				return downloaders.NewPastebinDownloaderWithParams(params...)
			},
		},

		transformerDict: map[string] func(params... []byte) (transformers.Transformer, error) {
			transformers.Base64TransformerName: func(params... []byte) (transformers.Transformer, error) {
				if len(params) != 0 {
					return nil, errors.New("base64 transformer doesn't accept any params")
				}

				return transformers.NewBase64Transformer(), nil
			},
			transformers.AESEncryptorName: func(params... []byte) (transformers.Transformer, error) {
				return encryptors.NewAESEncryptorWithParams(params...)
			},
		},

		uploaderDict: map[string] func(params... []byte) (uploaders.Uploader, error) {
			uploaders.PastebinUploaderName: func(params... []byte) (uploaders.Uploader, error) {
				return uploaders.NewPastebinUploaderWithParams(params...)
			},
		},
	}
}

func (br *builtinResolver) Resolve(urlPart string) (common.Nameable, error) {
	return nil, nil
}


func (br *builtinResolver) ResolveDownloader(name string, params... []byte) (downloaders.Downloader, error) {
	if downloaderFactory, ok := br.downloaderDict[name]; ok {
		return downloaderFactory(params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in downloader with name %s", name))
}

func (br *builtinResolver) ResolveTransformer(name string, params... []byte) (transformers.Transformer, error) {
	if transformerFactory, ok := br.transformerDict[name]; ok {
		return transformerFactory(params...)
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in transformer with name %s", name))
}

func (br *builtinResolver) ResolveUploader(name string, params... []byte) (uploaders.Uploader, error) {
	if uploaderFactory, ok := br.uploaderDict[name]; ok {
		return uploaderFactory(params...)
	}

	return nil, errors. New(fmt.Sprintf("there is no built-in uploader with name %s", name))
}