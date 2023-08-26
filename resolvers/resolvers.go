package resolvers

import (
	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/encryptors"
	"github.com/takahawk/shadownet/uploaders"
)

type Resolver interface {
	// urlPart - base64 decoded single part of shadownet URL, should be in format:
	//           Type:Name:Parameters
	Resolve(urlPart string) (common.Nameable, error)
	ResolveDownloader(name string, params... string) (downloaders.Downloader, error)
	ResolveTransformer(name string, params... string) (transformers.Transformer, error)
	ResolveEncryptor(name string, params... string) (encryptors.Encryptor, error)
	ResolveUploader(name string, params... string) (uploaders.Uploader, error)
}

// TODO: add plugin and/or socket and/or IPC bridge implementations