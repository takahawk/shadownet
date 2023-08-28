package resolvers

import (
	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type Resolver interface {
	// urlPart - base64 decoded single part of shadownet URL, should be in format:
	//           Type:Name:Parameters
	Resolve(urlPart string) (common.Component, error)
	ResolveDownloader(name string, params... []byte) (downloaders.Downloader, error)
	ResolveTransformer(name string, params... []byte) (transformers.Transformer, error)
	ResolveUploader(name string, params... []byte) (uploaders.Uploader, error)
}

// TODO: add plugin and/or socket and/or IPC bridge implementations