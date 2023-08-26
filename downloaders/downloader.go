package downloaders

import (
	"github.com/takahawk/shadownet/common"
)

type Downloader interface {
	common.Nameable
	Download(id string) (string, error)
}
