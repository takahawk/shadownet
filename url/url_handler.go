package url

import (
	"github.com/takahawk/shadownet/common"
)

const DownloaderURLPrefix = "down"
const TransformerURLPrefix = "trans"

// Tool to handle ShadowNet URLs
type UrlHandler interface {
	// Use ID and upload components to make ShadowNet URL
	MakeURL(id string, components... common.Component) (string, error)
	// Parse URL to get download components
	GetDownloadComponents(url string) ([]common.Component, error)
}