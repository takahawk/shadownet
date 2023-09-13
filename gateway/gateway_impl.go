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
	"github.com/takahawk/shadownet/storages"
)

type shadowGateway struct {
	logger  logger.Logger
	storage storages.Storage
	// TODO: cache pipelines?
}

type pipelineSpec struct {
	Components []struct {
		Name            string
		Params          []string
		IsParamsBase64d bool
	}
}

func NewShadowGateway(logger logger.Logger, storage storages.Storage) ShadownetGateway {
	// TODO: check for nil parameters
	return &shadowGateway{
		logger:  logger,
		storage: storage,
	}
}

func (sg *shadowGateway) Start(port int) error {
	r := mux.NewRouter()
	r.HandleFunc("/{shadowUrl}", sg.handleGatewayRequest).Methods("GET")
	r.HandleFunc("/setupPipeline/{pipelineName}", sg.handleSetupPipelineRequest).Methods("POST")
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
	vars := mux.Vars(req)
	pipelineName := vars["pipelineName"]
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}
	pipelineJson := string(b)
	_, err = sg.parsePipeline(pipelineJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sg.storage.SavePipelineJSON(pipelineName, pipelineJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: check if name exists
	sg.logger.Infof("Pipeline with name \"%s\" successfully added\n", pipelineName)
	fmt.Fprintf(w, "Pipeline with name \"%s\" successfully added\n", pipelineName)
}

func (sg *shadowGateway) handleUploadFileRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pipelineName := vars["pipelineName"]
	pipelineJson, err := sg.storage.LoadPipelineJSON(pipelineName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	pipeline, err := sg.parsePipeline(pipelineJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
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
	sg.logger.Infof("New page successfully uploaded: %s", url)
	w.Write(responseJson)
}

func (sg *shadowGateway) parsePipeline(pipelineJson string) (pipelines.UploadPipeline, error) {
	var request pipelineSpec
	err := json.Unmarshal([]byte(pipelineJson), &request)
	if err != nil {
		sg.logger.Errorf("%+v", err)
		return nil, err
	}
	if len(request.Components) == 0 {
		sg.logger.Error("Setup pipeline request with empty pipeline")
		return nil, err
	}

	var byteParams [][][]byte
	for j, component := range request.Components {
		byteParams = append(byteParams, nil)
		for _, param := range component.Params {
			if component.IsParamsBase64d {
				decoded, err := base64.StdEncoding.DecodeString(param)
				if err != nil {
					sg.logger.Errorf("%+v", err)
					return nil, err
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
			sg.logger.Errorf("%+v", err)
			return nil, err
		}
		err = pipeline.AddSteps(transformer)
		if err != nil {
			sg.logger.Errorf("%+v", err)
			return nil, err
		}
	}

	// TODO: mb double-check for uploader and send human-friendly error
	uploaderSpec := request.Components[len(request.Components)-1]

	uploader, err := resolver.ResolveUploader(uploaderSpec.Name, byteParams[len(request.Components)-1]...)
	if err != nil {
		sg.logger.Errorf("%+v", err)
		return nil, err
	}

	err = pipeline.AddSteps(uploader)
	if err != nil {
		sg.logger.Errorf("%+v", err)
		return nil, err
	}

	return pipeline, nil
}
