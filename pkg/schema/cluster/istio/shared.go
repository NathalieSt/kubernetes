package istio

type IstioPortProtocol string

const (
	TCP  IstioPortProtocol = "TCP"
	HTTP IstioPortProtocol = "HTTP"
)

type IstioPort struct {
	Number   int
	Name     string
	Protocol IstioPortProtocol
}
