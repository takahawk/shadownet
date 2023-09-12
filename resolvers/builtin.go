package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type builtinResolver struct {
	downloaderDict  map[string]func(params ...[]byte) (downloaders.Downloader, error)
	transformerDict map[string]func(params ...[]byte) (transformers.Transformer, error)
	uploaderDict    map[string]func(params ...[]byte) (uploaders.Uploader, error)
}

// NewBuiltinResolver returns new built-in resolver that is created to resolve
// components built-in into ShadowNet directly (TODO: provide list?)
func NewBuiltinResolver() Resolver {
	return &builtinResolver{
		downloaderDict: map[string]func(params ...[]byte) (downloaders.Downloader, error){
			downloaders.PastebinDownloaderName: func(params ...[]byte) (downloaders.Downloader, error) {
				return downloaders.NewPastebinDownloaderWithParams(params...)
			},
		},

		transformerDict: map[string]func(params ...[]byte) (transformers.Transformer, error){
			transformers.Base64TransformerName: func(params ...[]byte) (transformers.Transformer, error) {
				if len(params) != 0 {
					return nil, errors.New("base64 transformer doesn't accept any params")
				}

				return transformers.NewBase64Transformer(), nil
			},
			transformers.AESEncryptorName: func(params ...[]byte) (transformers.Transformer, error) {
				return transformers.NewAESEncryptorWithParams(params...)
			},
		},

		uploaderDict: map[string]func(params ...[]byte) (uploaders.Uploader, error){
			uploaders.PastebinUploaderName: func(params ...[]byte) (uploaders.Uploader, error) {
				return uploaders.NewPastebinUploaderWithParams(params...)
			},
		},
	}
}

// ResolveDownloader returns one of built-in downloaders
func (br *builtinResolver) ResolveDownloader(name string, params ...[]byte) (downloaders.Downloader, error) {
	if downloaderFactory, ok := br.downloaderDict[name]; ok {
		return downloaderFactory(params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in downloader with name %s", name))
}

// ResolveTransformer returns one of built-in transformers
func (br *builtinResolver) ResolveTransformer(name string, params ...[]byte) (transformers.Transformer, error) {
	if transformerFactory, ok := br.transformerDict[name]; ok {
		return transformerFactory(params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in transformer with name %s", name))
}

// ResolveUploader returns one of built-in uploaders
func (br *builtinResolver) ResolveUploader(name string, params ...[]byte) (uploaders.Uploader, error) {
	if uploaderFactory, ok := br.uploaderDict[name]; ok {
		return uploaderFactory(params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in uploader with name %s", name))
}
