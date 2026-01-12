package openwrt

import (
	"strings"
	"wnetctl/site"
)

func siteRequestToModel(request *site.SiteRequest) *SiteModel {
	model := new(SiteModel)
	model.SshKey = request.SshKey
	model.SshPublicKey = request.SshPublicKey
	model.Password = request.Password
	return model
}

func apRequestToModel(request *site.AccessPointRequest) *AccessPointModel {
	model := new(AccessPointModel)
	model.Model = request.Model
	model.Name = request.Name
	model.Ip = request.Ip
	return model
}

func ssidToSiteSsid(ssid *SSID) *site.SSID {
	sssid := new(site.SSID)
	sssid.Name = ssid.Name
	sssid.Auth = ssid.Auth
	sssid.Vlan = ssid.Vlan
	sssid.Password = ssid.Password
	return sssid
}

func siteSsidToSsid(sssid *site.SSID) *SSID {
	ssid := new(SSID)
	ssid.Name = sssid.Name
	ssid.Auth = sssid.Auth
	ssid.Vlan = sssid.Vlan
	ssid.Password = sssid.Password
	return ssid
}

func stationToSiteStation(station *Station) *site.Station {
	st := new(site.Station)
	st.Name = station.Name
	st.Mac = station.Mac
	st.Comment = station.Comment
	return st
}

func siteStationToStation(station *site.Station) *Station {
	st := new(Station)
	st.Name = station.Name
	st.Mac = strings.ToLower(station.Mac)
	st.Comment = station.Comment
	return st
}

func siteDeviceToDevice(stdev *site.AccessPointDevice) *AccessPointDevice {
	dev := new(AccessPointDevice)
	// TODO
	return dev
}

func deviceToSiteDevice(dev *AccessPointDevice) *site.AccessPointDevice {
	stdev := new(site.AccessPointDevice)
	// TODO
	return stdev
}

func wirelessAdapterToModel(adapter *WirelessAdapter) *WirelessAdapterModel {
	model := new(WirelessAdapterModel)
	model.Mac = adapter.Mac
	model.Channel = adapter.Channel
	model.Power = adapter.Power
	return model
}

func modelToWirelessAdapter(adapter *WirelessAdapter, model *WirelessAdapterModel, dev *DeviceWirelessAdapter) {
	adapter.Mac = model.Mac
	adapter.Channel = model.Channel
	adapter.Power = model.Power
	adapter.Device = new(DeviceWirelessAdapter)
	adapter.Device.Device = dev.Device
	adapter.Device.Interface = dev.Interface
	adapter.Device.Driver = dev.Driver
}
