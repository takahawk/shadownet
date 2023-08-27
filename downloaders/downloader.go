package downloaders

import (
	"github.com/takahawk/shadownet/common"
)

type Downloader interface {
	common.Component
	Download() (string, error)
}
