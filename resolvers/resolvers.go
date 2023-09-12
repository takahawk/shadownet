package resolvers

import (
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

// Resolver can be used to get ShadowNet component given a name, parameters and
// type of component
type Resolver interface {
	// ResolveDownloader returns Downloader by name and parameters
	ResolveDownloader(name string, params ...[]byte) (downloaders.Downloader, error)
	// ResolveDownloader returns Transformer by name and parameters
	ResolveTransformer(name string, params ...[]byte) (transformers.Transformer, error)
	// ResolveDownloader returns Uploader by name and parameters
	ResolveUploader(name string, params ...[]byte) (uploaders.Uploader, error)
}

// TODO: add plugin and/or socket and/or IPC bridge implementations
