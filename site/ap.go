package site

type AccessPoint interface {
	Configure() error
	AddNeighbour(neighbour AccessPoint) error
	RemoveNeighbour(neighbour AccessPoint) error
	AddSSID(ssid *SSID) error
	RemoveSSID(ssid *SSID) error
	AddStation(mac string) error
	RemoveStation(mac string) error
	Name() string
	ToResponse() *AccessPointResponse
}
