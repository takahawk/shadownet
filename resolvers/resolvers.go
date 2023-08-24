package resolvers

import (
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/encryptors"
)

type DownloaderResolver interface {
	ResolveDownloader(urlPart string) (downloaders.Downloader, error)
}

type TransformerResolver interface {
	ResolveTransformer(urlPart string) (transformers.Transformer, error)
}

type EncryptorResolver interface {
	ResolveEncryptor(urlPart string) (encryptors.Encryptor, error)
}

// TODO: add plugin and/or socket and/or IPC bridge implementations