package pipelines

import (
	"github.com/takahawk/shadownet/common"
)

const DownloaderURLPrefix = "down"
const EncryptorURLPrefix = "enc"
const TransformerURLPrefix = "trans"

type UploadPipeline interface {
	// key is only for encryptors, should be nil for other parts of pipeline
	AddSteps(components... common.Component) error
	Upload(data []byte) (url string, err error)
}