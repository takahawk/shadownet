package pipelines

import (
	"errors"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/uploaders"
	"github.com/takahawk/shadownet/url"
	"github.com/takahawk/shadownet/transformers"
)

type uploadPipeline struct {
	finalized bool
	urlHandler url.UrlHandler
	steps []common.Component
}

func NewUploadPipeline() UploadPipeline {
	return &uploadPipeline{}
}

func (up *uploadPipeline) AddSteps(components... common.Component) error {
	// there can be only one uploader and it will always be the last one
	if up.finalized {
		return errors.New("can't add another step. There is already an uploader")
	}

	for i, component := range components {
		if _, ok := component.(uploaders.Uploader); ok {
			if i != len(components) - 1 {
				return errors.New("uploader can be only the last step")
			} else {
				up.finalized = true
			}
		}
	}
	

	up.steps = append(up.steps, components...)

	return nil
}


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