package storages

import "github.com/takahawk/shadownet/models"

// Storage
type Storage interface {
	PipelineStorage
}

// PipelineStorage is used to persistently store pipelines in JSON form
type PipelineStorage interface {
	// ListPipelineSpecs returns slice of all pipeline specifications
	// that are exist in storage
	ListPipelineSpecs() ([]*models.PipelineSpec, error)
	// SavePipelineSpec stores pipeline specification
	SavePipelineSpec(spec *models.PipelineSpec) error
	// UpdatePipelineSpec overwrites pipeline specification with a given name
	UpdatePipelineSpec(spec *models.PipelineSpec) error
	// LoadPipelineSpec returns pipeline specification with a given name
	LoadPipelineSpec(name string) (*models.PipelineSpec, error)
	// DeletePipelineSpec removes pipeline specification with a given name from
	// storage
	DeletePipelineSpec(name string) error
}
