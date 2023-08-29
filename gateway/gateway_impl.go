package gateway

import (
	"fmt"
	"net/http"

	"github.com/takahawk/shadownet/pipelines"
)

type shadowGateway struct {}

func NewShadowGateway() ShadownetGateway {
	return &shadowGateway{}
}

func (sg *shadowGateway) Start(port int) error {
	http.HandleFunc("/", handleGatewayRequest)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleGatewayRequest(w http.ResponseWriter, req *http.Request) {
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