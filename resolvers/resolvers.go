package resolvers

import (
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type Resolver interface {
	ResolveDownloader(name string, params... []byte) (downloaders.Downloader, error)
	ResolveTransformer(name string, params... []byte) (transformers.Transformer, error)
	ResolveUploader(name string, params... []byte) (uploaders.Uploader, error)
}

// TODO: add plugin and/or socket and/or IPC bridge implementations