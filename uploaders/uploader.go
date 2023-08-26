package uploaders

import (
	"github.com/takahawk/shadownet/common"
)


type Uploader interface {
	common.Nameable
	// TODO: mb add some generic `params` to gain more control on specific upload?
	Upload(content []byte) (id string, err error)
}