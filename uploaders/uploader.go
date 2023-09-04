package uploaders

import (
	"github.com/takahawk/shadownet/common"
)

// Uploader is component used to upload data from a specific storage
type Uploader interface {
	common.Component
	// Upload uploads data in a byte array to storage returning id which can
	// be used by corresponding Downloader to get that data
	Upload(content []byte) (id string, err error)
}
