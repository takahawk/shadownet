package pipelines

import (
	"errors"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/downloaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
	"github.com/takahawk/shadownet/url"
)

type uploadPipeline struct {
	finalized  bool
	urlHandler url.UrlHandler
	steps      []common.Component
}

// NewUploadPipeline returns new empty upload pipeline
func NewUploadPipeline() UploadPipeline {
	return &uploadPipeline{
		urlHandler: url.NewUrlHandler(),
	}
}

// AddSteps adds components to upload pipeline
func (up *uploadPipeline) AddSteps(components ...common.Component) error {
	// there can be only one uploader and it will always be the last one
	if up.finalized {
		return errors.New("can't add another step. There is already an uploader")
	}

	for i, component := range components {
		if _, ok := component.(uploaders.Uploader); ok {
			if i != len(components)-1 {
				return errors.New("uploader can be only the last step")
			} else {
				up.finalized = true
			}
		}
	}

	up.steps = append(up.steps, components...)

	return nil
}

// Upload runs the whole pipeline to upload given data returning ShadowNet URL
func (up *uploadPipeline) Upload(data []byte) (url string, err error) {
	if !up.finalized {
		return "", errors.New("pipeline should be finalized with uploader")
	}

	var id string
	for _, step := range up.steps {
		switch step := step.(type) {
		case transformers.Transformer:
			data, err = step.ForwardTransform(data)
		case uploaders.Uploader:
			id, err = step.Upload(data)
		}

		if err != nil {
			return "", err
		}
	}

	return up.urlHandler.MakeURL(id, up.steps...)
}

type downloadPipeline struct {
	steps []common.Component
}

// NewDownloadPipeline returns new empty download pipeline
func NewDownloadPipeline() DownloadPipeline {
	return &downloadPipeline{}
}

// NewDownloadPipeline constructs from ShadowNet URL new download pipeline
// that can be used to download and decode data from it
func NewDownloadPipelineByURL(shadowUrl string) (DownloadPipeline, error) {
	urlHandler := url.NewUrlHandler()
	pipeline := NewDownloadPipeline()
	components, err := urlHandler.GetDownloadComponents(shadowUrl)
	if err != nil {
		return nil, err
	}
	pipeline.AddSteps(components...)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// AddSteps adds components to download pipeline
func (dp *downloadPipeline) AddSteps(components ...common.Component) error {
	// there can be only one downloader and it will always be the first one
	if len(dp.steps) == 0 && len(components) != 0 {
		if _, ok := components[0].(downloaders.Downloader); !ok {
			return errors.New("first step should be downloader")
		}
	}

	for i, component := range components {
		if _, ok := component.(downloaders.Downloader); ok {
			if len(dp.steps) != 0 && i != 0 {
				return errors.New("downloader can be only the first step")
			}
		}
	}

	dp.steps = append(dp.steps, components...)

	return nil
}

// Download runs the whole pipeline to retrieve and trasform data returning the
// final result
func (dp *downloadPipeline) Download() (data []byte, err error) {
	if len(dp.steps) == 0 {
		return nil, errors.New("empty download pipeline")
	}
	for _, component := range dp.steps {
		switch component := component.(type) {
		case downloaders.Downloader:
			data, err = component.Download()
		case transformers.Transformer:
			data, err = component.ReverseTransform(data)
		}

		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
