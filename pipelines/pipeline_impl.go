package pipelines

import (
	"encoding/base64"

	"errors"
	"strings"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/uploaders"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/encryptors"
)

// TODO: mb make public?
type step struct {
	nameable common.Nameable
	// TODO: mb move params to nameable itself?
	params [][]byte
}
type uploadPipeline struct {
	finalized bool
	steps []step
}

func NewUploadPipeline() UploadPipeline {
	return &uploadPipeline{}
}

func (up *uploadPipeline) AddStep(nameable common.Nameable, params... []byte) error {
	// there can be only one uploader and it will always be the last one
	if up.finalized {
		return errors.New("can't add another step. There is already an uploader")
	}

	up.steps = append(up.steps, step{ nameable: nameable, params: params })
	if _, ok := nameable.(uploaders.Uploader); ok {
		up.finalized = true
	}

	return nil
}


func (up *uploadPipeline) Upload(data []byte) (url string, err error) {
	if !up.finalized {
		return "", errors.New("pipeline should be finalized with uploader")
	}

	var urlParts []string
	for _, step := range up.steps {
		switch nameable := step.nameable.(type) {
		case encryptors.Encryptor:
			data, err = nameable.Encrypt(data)
		case transformers.Transformer:
			data, err = nameable.ForwardTransform(data)
		case uploaders.Uploader:
			var id string
			id, err = nameable.Upload(data)
			step.params = [][]byte{[]byte(id)}
		}

		if err != nil {
			return "", err
		}

		urlParts = append(urlParts, getURLPart(step.nameable, step.params...))
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
func getURLPart(nameable common.Nameable, params... []byte) string {
	var sb strings.Builder

	switch nameable.(type) {
	case encryptors.Encryptor:
		sb.WriteString(EncryptorURLPrefix)
	case transformers.Transformer:
		sb.WriteString(TransformerURLPrefix)
	case uploaders.Uploader:
		sb.WriteString(DownloaderURLPrefix)
	}
	sb.WriteString("_")
	sb.WriteString(nameable.Name())
	sb.WriteString(":")
	for _, param := range params {
		sb.WriteString(base64.StdEncoding.EncodeToString(param))
		sb.WriteString(",")
	}
	return base64.StdEncoding.EncodeToString([]byte(sb.String()))
}