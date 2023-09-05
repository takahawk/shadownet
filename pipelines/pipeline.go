package pipelines

import (
	"github.com/takahawk/shadownet/common"
)

// UploadPipeline groups components to transform and then upload data to
// storage.
type UploadPipeline interface {
	// AddSteps add components to upload pipeline. The last component should
	// always be uploader and uploader can be only the last component in
	// pipeline
	AddSteps(components ...common.Component) error
	// Upload runs the whole pipeline for a given data, transforming and then
	// uploading it to storage.
	// ShadowNet URL to access the data is returned afterwards
	Upload(data []byte) (url string, err error)
}

// DownloadPipeline groups components to transform and then upload data to
// storage.
type DownloadPipeline interface {
	// AddSteps add components to download pipeline. The first component should
	// always be download and it can be only the first component in
	// pipeline
	AddSteps(components ...common.Component) error
	// Download run the whole pipeline downloading data from storage and then
	// transforming it and returning back. It is supposed to do effectively
	// the reverse process of what Uploader.Upload have done
	Download() (data []byte, err error)
}
