package url

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/resolvers"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type urlHandler struct {
	resolver resolvers.Resolver
}

func NewUrlHandler() UrlHandler {
	resolver := resolvers.NewBuiltinResolver()
	return &urlHandler{
		resolver: resolver,
	}
}


func (uh *urlHandler) MakeURL(id string, components... common.Component) (string, error) {
	var urlParts []string

	for _, component := range components {
		switch component := component.(type) {
		case transformers.Transformer:
			urlParts = append(urlParts, getURLPart(component, component.Params()...))
		case uploaders.Uploader:
			// mb double-check for uploader to be only the last component?
			urlParts = append(urlParts, getURLPart(component, []byte(id)))
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

func (uh *urlHandler) GetDownloadComponents(url string) ([]common.Component, error) {
	// TODO: impl
	return nil, errors.New("Not implemented")
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