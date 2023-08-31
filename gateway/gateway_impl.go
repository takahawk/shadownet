package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/takahawk/shadownet/pipelines"
	"github.com/takahawk/shadownet/resolvers"
)

type shadowGateway struct {
	// TODO: store in db?
	pipelines map[string] pipelines.UploadPipeline
}

type setupPipelineRequest struct {
	Name string
	Components []struct {
		Name string
		Params []string
		IsParamsBase64d bool
	}
}

func NewShadowGateway() ShadownetGateway {
	return &shadowGateway{}
}

func (sg *shadowGateway) Start(port int) error {
	http.HandleFunc("/", sg.handleGatewayRequest)
	http.HandleFunc("/setupPipeline", sg.handleSetupPipelineRequest)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (sg *shadowGateway) handleGatewayRequest(w http.ResponseWriter, req *http.Request) {
	// TODO: handle empty case separately
	shadowUrl := req.URL.Path[1:]
	pipeline, err := pipelines.NewDownloadPipelineByURL(shadowUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", err)
		return
	}
	data, err := pipeline.Download()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", err)
	}
	
	fmt.Fprintf(w, string(data))
}

func (sg *shadowGateway) handleSetupPipelineRequest(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var request setupPipelineRequest
		err = json.Unmarshal(b, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(request.Components) == 0 {
			http.Error(w, "there should be at least one component", http.StatusBadRequest)
			return
		}

		var byteParams [][][]byte
		for j, component := range request.Components {
			byteParams = append(byteParams, nil)
			for _, param := range component.Params {
				if component.IsParamsBase64d {
					decoded, err := base64.StdEncoding.DecodeString(param)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					byteParams[j] = append(byteParams[j], decoded)
				} else {
					byteParams[j] = append(byteParams[j], []byte(param))
				}
			} 
		}

		pipeline := pipelines.NewUploadPipeline()

		resolver := resolvers.NewBuiltinResolver()
		uploaderSpec := request.Components[0]

		uploader, err := resolver.ResolveUploader(uploaderSpec.Name, byteParams[0]...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		err = pipeline.AddSteps(uploader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i := 1; i < len(request.Components); i++ {
			
			transformer, err := resolver.ResolveTransformer(request.Components[i].Name, byteParams[i]...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = pipeline.AddSteps(transformer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		sg.pipelines[request.Name] = pipeline
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "method %s not supported", req.Method)
	}

	fmt.Printf("%+v", sg.pipelines)
}

func (sg *shadowGateway) handleUploadFileRequest(w http.ResponseWriter, req *http.Request) {
	// TODO: impl
}