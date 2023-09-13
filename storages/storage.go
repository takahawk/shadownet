package storages

type PipelineJSONsListEntry struct {
	Name string
	JSON string
}

type Storage interface {
	PipelineStorage
}

type PipelineStorage interface {
	ListPipelineJSONs() ([]PipelineJSONsListEntry, error)
	SavePipelineJSON(name string, json string) error
	UpdatePipelineJSON(name string, json string) error
	LoadPipelineJSON(name string) (json string, err error)
	DeletePipelineJSON(name string) error
}
