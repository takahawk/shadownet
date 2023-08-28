package pipelines

import (
	"encoding/base64"

	"errors"
	"strings"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/uploaders"
	"github.com/takahawk/shadownet/transformers"
)

type uploadPipeline struct {
	finalized bool
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

	var urlParts []string
	for _, step := range up.steps {
		switch step := step.(type) {
		case transformers.Transformer:
			data, err = step.ForwardTransform(data)
			urlParts = append(urlParts, getURLPart(step, step.Params()...))
		case uploaders.Uploader:
			var id string
			id, err = step.Upload(data)
			urlParts = append(urlParts, getURLPart(step, []byte(id)))
		}

		if err != nil {
			return "", err
		}

		
	}

	var sb strings.Builder
	for i := len(urlParts) - 1; i > 0; i-- {
		sb.WriteString(urlParts[i])
		sb.WriteString(".")
	}
	sb.WriteString(urlParts[0])

	return sb.String(), nil
	
}

// [Type]_[ID]:[Base64dCommaSeparatedParameters]
func getURLPart(component common.Component, params... []byte) string {
	var sb strings.Builder

	switch component.(type) {
	case transformers.Transformer:
		sb.WriteString(TransformerURLPrefix)
	case uploaders.Uploader:
		sb.WriteString(DownloaderURLPrefix)
	}
	sb.WriteString("_")
	sb.WriteString(component.Name())
	sb.WriteString(":")
	for _, param := range params {
		sb.WriteString(base64.StdEncoding.EncodeToString(param))
		sb.WriteString(",")
	}
	return base64.StdEncoding.EncodeToString([]byte(sb.String()))
}