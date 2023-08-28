package pipelines

import (
	"github.com/takahawk/shadownet/common"
)

type UploadPipeline interface {
	AddSteps(components... common.Component) error
	Upload(data []byte) (url string, err error)
}
