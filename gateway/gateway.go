package gateway

type ShadownetGateway interface {
	Start(port int) error
}