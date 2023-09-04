package downloaders

import (
	"github.com/takahawk/shadownet/common"
)

// Downloader is component used to download data from a specific storage
type Downloader interface {
	common.Component
	// Download downloads data to a byte array
	Download() ([]byte, error)
}
