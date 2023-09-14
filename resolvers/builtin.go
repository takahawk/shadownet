package resolvers

import (
	"errors"
	"fmt"

	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type builtinResolver struct {
	logger          logger.Logger
	downloaderDict  map[string]func(logger logger.Logger, params ...[]byte) (downloaders.Downloader, error)
	transformerDict map[string]func(logger logger.Logger, params ...[]byte) (transformers.Transformer, error)
	uploaderDict    map[string]func(logger logger.Logger, params ...[]byte) (uploaders.Uploader, error)
}

// NewBuiltinResolver returns new built-in resolver that is created to resolve
// components built-in into ShadowNet directly (TODO: provide list?)
func NewBuiltinResolver(log logger.Logger) Resolver {
	return &builtinResolver{
		logger: log,
		downloaderDict: map[string]func(logger logger.Logger, params ...[]byte) (downloaders.Downloader, error){
			downloaders.PastebinDownloaderName: downloaders.NewPastebinDownloaderWithParams,
		},

		transformerDict: map[string]func(logger logger.Logger, params ...[]byte) (transformers.Transformer, error){
			transformers.Base64TransformerName: transformers.NewBase64TransformerWithParams,
			transformers.AESEncryptorName:      transformers.NewAESEncryptorWithParams,
		},

		uploaderDict: map[string]func(logger logger.Logger, params ...[]byte) (uploaders.Uploader, error){
			uploaders.PastebinUploaderName: uploaders.NewPastebinUploaderWithParams,
			uploaders.DropboxUploaderName:  uploaders.NewDropboxUploaderWithParams,
		},
	}
}

// ResolveDownloader returns one of built-in downloaders
func (br *builtinResolver) ResolveDownloader(name string, params ...[]byte) (downloaders.Downloader, error) {
	if downloaderFactory, ok := br.downloaderDict[name]; ok {
		return downloaderFactory(br.logger, params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in downloader with name %s", name))
}

// ResolveTransformer returns one of built-in transformers
func (br *builtinResolver) ResolveTransformer(name string, params ...[]byte) (transformers.Transformer, error) {
	if transformerFactory, ok := br.transformerDict[name]; ok {
		return transformerFactory(br.logger, params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in transformer with name %s", name))
}

// ResolveUploader returns one of built-in uploaders
func (br *builtinResolver) ResolveUploader(name string, params ...[]byte) (uploaders.Uploader, error) {
	if uploaderFactory, ok := br.uploaderDict[name]; ok {
		return uploaderFactory(br.logger, params...)
	}

	return nil, errors.New(fmt.Sprintf("there is no built-in uploader with name %s", name))
}
