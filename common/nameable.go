package common

// Component is the main building block of ShadowNet that does something with
// a data (downloads/uploads, transforms it any way etc.)
type Component interface {
	// Name returns the name of component
	Name() string
	// Params is the parameters which are supposed to be used by the specific
	// instance of component
	Params() [][]byte
}
