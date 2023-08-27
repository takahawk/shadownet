package common

type Component interface {
	Name() string
	Params() [][]byte
}