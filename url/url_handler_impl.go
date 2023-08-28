package url

import (
	"fmt"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"github.com/takahawk/shadownet/common"
	"github.com/takahawk/shadownet/resolvers"
	"github.com/takahawk/shadownet/transformers"
	"github.com/takahawk/shadownet/uploaders"
)

type urlHandler struct {
	resolver resolvers.Resolver
}

var urlPartPattern = regexp.MustCompile(`(trans|down)_(.+):(.*)`)

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
	components := make([]common.Component, 0)
	for _, urlPart := range strings.Split(url, ".") {
		urlPart, err := base64.StdEncoding.DecodeString(urlPart)
		if err != nil {
			return nil, err
		}

		groups := urlPartPattern.FindStringSubmatch(string(urlPart))
		if len(groups) != 4 {
			return nil, errors.New(fmt.Sprintf("invalid url part: %s", urlPart))
		}
		prefix := groups[1]
		name := groups[2]
		strParams := groups[3]

		var params [][]byte
		if strParams != "" {
			for _, strParam := range strings.Split(strParams, ",") {
				param, err := base64.StdEncoding.DecodeString(strParam)
				if err != nil {
					return nil, err
				}
				params = append(params, param)
			}
		}

		var component common.Component
		switch prefix {
		case DownloaderURLPrefix:
			component, err = uh.resolver.ResolveDownloader(name, params...)
		case TransformerURLPrefix:
			component, err = uh.resolver.ResolveTransformer(name, params...)
		}

		if err != nil {
			return nil, err
		}

		components = append(components, component)

	}
	return components, nil
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
	for i, param := range params {
		sb.WriteString(base64.StdEncoding.EncodeToString(param))
		if i != len(params) - 1 {
			sb.WriteString(",")
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(sb.String()))
}