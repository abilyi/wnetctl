package site

import "io"

type SiteManager interface {
	GetSite() *SiteResponse
	AddAccessPoint(model *AccessPointRequest) (AccessPoint, error)
	GetAccessPoints() []*AccessPointResponse
	//UpdateAccessPoint(*AccessPoint) error
	RemoveAccessPoint(name string) error
	AddSSID(*SSID) error
	GetSSIDs() []*SSID
	UpdateSSID(*SSID) error
	RemoveSSID(string) error
	AddStation(string, *Station) error
	GetStations(ssid string) ([]*Station, error)
	RemoveStation(ssidName, mac string) error
	AddDeviceType(device *AccessPointDevice) error
	RemoveDeviceType(deviceType string) error
	GetDeviceTypes() []*AccessPointDevice
	Export(dest io.Writer) error
}
