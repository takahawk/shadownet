package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/pipelines"
	"github.com/takahawk/shadownet/resolvers"
)

type shadowGateway struct {
	// TODO: store in db?
	logger    logger.Logger
	pipelines map[string]pipelines.UploadPipeline
}

type setupPipelineRequest struct {
	Name       string
	Components []struct {
		Name            string
		Params          []string
		IsParamsBase64d bool
	}
}

func NewShadowGateway(logger logger.Logger) ShadownetGateway {
	return &shadowGateway{
		logger:    logger,
		pipelines: make(map[string]pipelines.UploadPipeline),
	}
}

func (sg *shadowGateway) Start(port int) error {
	r := mux.NewRouter()
	r.HandleFunc("/{shadowUrl}", sg.handleGatewayRequest).Methods("GET")
	r.HandleFunc("/setupPipeline", sg.handleSetupPipelineRequest).Methods("POST")
	r.HandleFunc("/upload/{pipelineName}", sg.handleUploadFileRequest).Methods("POST")
	http.Handle("/", r)
	sg.logger.Infof("Starting ShadowNet gateway on port %d", port)

	// TODO: mb add http access/error logs?
	// https://stackoverflow.com/questions/20987752/how-to-setup-access-error-log-for-http-listenandserve
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (sg *shadowGateway) handleGatewayRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	shadowUrl := vars["shadowUrl"]
	pipeline, err := pipelines.NewDownloadPipelineByURL(shadowUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		fmt.Fprintf(w, "%+v", err)
		return
	}
	data, err := pipeline.Download()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	fmt.Fprintf(w, string(data))
}

func (sg *shadowGateway) handleSetupPipelineRequest(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var request setupPipelineRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sg.logger.Errorf("%+v", err)
		return
	}
	if len(request.Components) == 0 {
		http.Error(w, "there should be at least one component", http.StatusBadRequest)
		sg.logger.Error("Setup pipeline request with empty pipeline")
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
					sg.logger.Errorf("%+v", err)
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

	for i := 0; i < len(request.Components)-1; i++ {
		transformer, err := resolver.ResolveTransformer(request.Components[i].Name, byteParams[i]...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			sg.logger.Errorf("%+v", err)
			return
		}
		err = pipeline.AddSteps(transformer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			sg.logger.Errorf("%+v", err)
			return
		}
	}

	// TODO: mb double-check for uploader and send human-friendly error
	uploaderSpec := request.Components[len(request.Components)-1]

	uploader, err := resolver.ResolveUploader(uploaderSpec.Name, byteParams[len(request.Components)-1]...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sg.logger.Errorf("%+v", err)
		return
	}

	err = pipeline.AddSteps(uploader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		sg.logger.Errorf("%+v", err)
		return
	}

	// TODO: check if name exists
	sg.pipelines[request.Name] = pipeline

	fmt.Printf("Pipeline with name \"%s\" successfully added\n", request.Name)
	fmt.Fprintf(w, "Pipeline with name \"%s\" successfully added\n", request.Name)
}

func (sg *shadowGateway) handleUploadFileRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pipelineName := vars["pipelineName"]

	if _, ok := sg.pipelines[pipelineName]; !ok {
		http.Error(w, fmt.Sprintf("there is no pipeline with name %s", pipelineName), http.StatusNotFound)
		sg.logger.Errorf("Pipeline with name \"%s\" is not found", pipelineName)
		return
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}
	defer file.Close()

	// TODO: mb better to make it work with Reader interface
	b, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}

	pipeline := sg.pipelines[pipelineName]
	url, err := pipeline.Upload(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}

	response := make(map[string]string)
	response["url"] = url
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}
	w.Write(responseJson)
}
