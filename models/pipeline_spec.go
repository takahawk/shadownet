package models

// PipelineSpec is specification of pipeline with name and its components
type PipelineSpec struct {
	Name       string `json:"name"`
	Components []struct {
		Name            string   `json:"name"`
		Params          []string `json:"params"`
		IsParamsBase64d bool     `json:"isParamsBase64d"`
	} `json:"components"`
}
