package storages

type Storage interface {
	PipelineStorage
}

type PipelineStorage interface {
	SavePipelineJSON(name string, json string) error
	LoadPipelineJSON(name string) (json string, err error)
}
