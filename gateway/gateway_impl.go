package gateway

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/models"
	"github.com/takahawk/shadownet/pipelines"
	"github.com/takahawk/shadownet/resolvers"
	"github.com/takahawk/shadownet/storages"
)

type shadowGateway struct {
	logger  logger.Logger
	storage storages.Storage
	// TODO: cache pipelines?
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
	r.HandleFunc("/pipelines", sg.handleListPipelinesRequest).Methods(http.MethodGet)
	r.HandleFunc("/{shadowUrl}", sg.handleGatewayRequest).Methods(http.MethodGet)
	r.HandleFunc("/pipelines", sg.handleAddPipelineRequest).Methods(http.MethodPost)
	r.HandleFunc("/pipelines", sg.handleUpdatePipelineRequest).Methods(http.MethodPut)
	r.HandleFunc("/pipelines/{pipelineName}", sg.handleDeletePipelineRequest).Methods(http.MethodDelete)

	r.HandleFunc("/pipelines/{pipelineName}/upload", sg.handleUploadFileRequest).Methods(http.MethodPost)
	http.Handle("/", r)
	sg.logger.Infof("Starting ShadowNet gateway on port %d", port)

	// TODO: mb add http access/error logs?
	// https://stackoverflow.com/questions/20987752/how-to-setup-access-error-log-for-http-listenandserve
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (sg *shadowGateway) handleGatewayRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	shadowUrl := vars["shadowUrl"]
	pipeline, err := pipelines.NewDownloadPipelineByURL(sg.logger, shadowUrl)
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

func (sg *shadowGateway) handleListPipelinesRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(w)
	pipelineSpecs, err := sg.storage.ListPipelineSpecs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}

	data, err := json.Marshal(pipelineSpecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}

	fmt.Fprintf(w, string(data))
}

func (sg *shadowGateway) handleAddPipelineRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(w)
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}
	pipelineJson := string(b)
	var pipelineSpec models.PipelineSpec
	err = json.Unmarshal([]byte(pipelineJson), &pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("Error unmarshaling pipeline spec: %+v", err)
		return
	}
	_, err = sg.parsePipeline(&pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("Error parsing pipeline: %+v", err)
		return
	}

	err = sg.storage.SavePipelineSpec(&pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sg.logger.Infof("Pipeline with name \"%s\" successfully added\n", pipelineSpec.Name)
	fmt.Fprintf(w, "Pipeline with name \"%s\" successfully added\n", pipelineSpec.Name)
}

func (sg *shadowGateway) handleUpdatePipelineRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(w)

	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("%+v", err)
		return
	}
	pipelineJson := string(b)
	var pipelineSpec models.PipelineSpec
	err = json.Unmarshal([]byte(pipelineJson), &pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sg.logger.Errorf("Error unmarshaling pipeline spec: %+v", err)
		return
	}
	_, err = sg.parsePipeline(&pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sg.storage.UpdatePipelineSpec(&pipelineSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sg.logger.Infof("Pipeline with name \"%s\" successfully updated\n", pipelineSpec.Name)
	fmt.Fprintf(w, "Pipeline with name \"%s\" successfully updated\n", pipelineSpec.Name)
}

func (sg *shadowGateway) handleDeletePipelineRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(w)
	vars := mux.Vars(req)
	pipelineName := vars["pipelineName"]

	err := sg.storage.DeletePipelineSpec(pipelineName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sg.logger.Infof("Pipeline with name \"%s\" successfully deleted\n", pipelineName)
	fmt.Fprintf(w, "Pipeline with name \"%s\" successfully deleted\n", pipelineName)
}

func (sg *shadowGateway) handleUploadFileRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pipelineName := vars["pipelineName"]
	pipelineSpec, err := sg.storage.LoadPipelineSpec(pipelineName)

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

	pipeline, err := sg.parsePipeline(pipelineSpec)
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

func (sg *shadowGateway) parsePipeline(pipelineSpec *models.PipelineSpec) (pipelines.UploadPipeline, error) {
	if len(pipelineSpec.Components) == 0 {
		sg.logger.Error("Setup pipeline request with empty pipeline")
		return nil, errors.New("setup pipeline request with empty pipeline")
	}

	var byteParams [][][]byte
	for j, component := range pipelineSpec.Components {
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

	pipeline := pipelines.NewUploadPipeline(sg.logger)

	resolver := resolvers.NewBuiltinResolver(sg.logger)

	for i := 0; i < len(pipelineSpec.Components)-1; i++ {
		transformer, err := resolver.ResolveTransformer(pipelineSpec.Components[i].Name, byteParams[i]...)
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
	uploaderSpec := pipelineSpec.Components[len(pipelineSpec.Components)-1]

	uploader, err := resolver.ResolveUploader(uploaderSpec.Name, byteParams[len(pipelineSpec.Components)-1]...)
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

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
