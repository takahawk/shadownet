package storages

// PipelineJSONsListEntry combines pipeline JSON with its name
type PipelineJSONsListEntry struct {
	Name string
	JSON string
}

// Storage
type Storage interface {
	PipelineStorage
}

// PipelineStorage is used to persistently store pipelines in JSON form
type PipelineStorage interface {
	// ListPipelineJSONs returns slice of all pipeline JSONs (with names)
	// that are exist in storage
	ListPipelineJSONs() ([]PipelineJSONsListEntry, error)
	// SavePipelineJSON stores pipeline JSON with a given name
	SavePipelineJSON(name string, json string) error
	// UpdatePipelineJSON overwrites pipeline JSON with a given name
	UpdatePipelineJSON(name string, json string) error
	// LoadPipelineJSON returns pipeline JSON with a given name
	LoadPipelineJSON(name string) (json string, err error)
	// DeletePipelineJSON removes pipeline JSON with a given name from storage
	DeletePipelineJSON(name string) error
}
