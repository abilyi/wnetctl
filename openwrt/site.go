package openwrt

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"slices"
	"strings"
	"wnetctl/site"
	"wnetctl/util"
)

const pluginName = "openwrt"

type Station struct {
	Name    string
	Mac     string
	Comment string
}

type SSID struct {
	Name        string
	Auth        string
	Password    string
	Vlan        int
	Restricted  bool
	Whitelisted bool
	Stations    []*Station
}

type Site struct {
	path         string
	name         string
	plugin       string
	sshKey       string
	sshPublicKey string
	password     string
	country      string
	suffix2      string
	suffix5      string
	accessPoints map[string]*AccessPoint
	ssids        []*SSID
	devices      map[string]*AccessPointDevice
}

type SiteModel struct {
	Plugin       string
	SshKey       string
	SshPublicKey string
	Password     string
	AccessPoints []*AccessPointModel
	Ssids        []*SSID
	Devices      []*AccessPointDevice
}

func NewSiteManager(name, path string) (site.SiteManager, error) {
	ste := new(Site)
	if err := util.ReadObject(path, ste); err != nil {
		return nil, err
	}
	ste.name = name
	ste.path = path
	ste.plugin = pluginName
	if err := ste.load(); err != nil {
		return nil, err
	}
	return ste, nil
}

func CreateSiteManager(name, path string, request *site.SiteRequest) (site.SiteManager, error) {
	ste := new(Site)
	ste.plugin = pluginName
	ste.name = name
	ste.path = path
	if err := ste.init(siteRequestToModel(request)); err != nil {
		return nil, err
	} else {
		if err = ste.save(); err != nil {
			return nil, err
		} else {
			return ste, nil
		}
	}
}

func (this *Site) GetSite() *site.SiteResponse {
	response := new(site.SiteResponse)
	response.Password = this.password
	response.SshKey = this.sshKey             // FIXME load key content instead
	response.SshPublicKey = this.sshPublicKey // FIXME load public key content instead
	response.Country = this.country
	response.SsidSuffix2 = this.suffix2
	response.SsidSuffix5 = this.suffix5
	response.Devices = make([]*site.AccessPointDevice, len(this.devices))
	i := 0
	for _, device := range this.devices {
		response.Devices[i] = new(site.AccessPointDevice)
		response.Devices[i].Name = device.Name
		response.Devices[i].Model = device.Model
		response.Devices[i].Cpu = device.Cpu
		response.Devices[i].Architecture = device.Architecture
		response.Devices[i].BridgedWiredDevice = device.BridgedWiredDevice
		response.Devices[i].WLan2.Device = device.Wlan2.Device
	}

	return response
}

func (this *Site) AddAccessPoint(request *site.AccessPointRequest) (site.AccessPoint, error) {
	dev, ok := this.devices[request.Model]
	if !ok || dev == nil {
		return nil, errors.New("Unknown device type \"" + request.Model + "\"")
	}

	accessPoint, err := CreateAccessPoint(request, this)
	if err == nil {
		err = accessPoint.Configure()
	}
	if err != nil {
		return nil, err
	}
	processed := make([]string, 0, len(this.accessPoints))
	for _, ap := range this.accessPoints {
		err := ap.AddNeighbour(accessPoint)
		if err != nil {
			for _, apName := range processed {
				this.accessPoints[apName].RemoveNeighbour(accessPoint)
			}
			return nil, err
		}
	}
	this.accessPoints[accessPoint.Name()] = accessPoint
	if err := this.save(); err != nil {
		return nil, err
	}
	return accessPoint, nil
}

func (this *Site) GetAccessPoints() []*site.AccessPointResponse {
	aps := make([]*site.AccessPointResponse, len(this.accessPoints))
	ix := 0
	for _, ap := range this.accessPoints {
		aps[ix] = ap.ToResponse()
		ix++
	}
	return aps
}

func (this *Site) RemoveAccessPoint(name string) error {
	accessPoint, ok := this.accessPoints[name]
	if !ok {
		return errors.New("Unknown access point \"" + name + "\"")
	}
	//var ix int
	//for ix = 0; ix < len(this.accessPoints) && this.accessPoints[ix].Name != name; ix++ {
	//}
	//if ix == len(this.accessPoints) {
	//	return errors.New("access point not found")
	//}
	for nm, ap := range this.accessPoints {
		if nm != name {
			ap.RemoveNeighbour(accessPoint)
		}
	}
	delete(this.accessPoints, name)
	return this.save()
}

func (this *Site) AddSSID(ssid *site.SSID) error {
	processed := make([]string, 0, len(this.accessPoints))
	for name, ap := range this.accessPoints {
		if err := ap.AddSSID(ssid); err != nil {
			for _, nm := range processed {
				this.accessPoints[nm].RemoveSSID(ssid)
			}
			return err
		}
		processed = append(processed, name)
	}
	this.ssids = append(this.ssids, siteSsidToSsid(ssid))
	return this.save()
}

func (this *Site) GetSSIDs() []*site.SSID {
	ssids := make([]*site.SSID, len(this.ssids))
	for i, ssid := range this.ssids {
		ssids[i] = ssidToSiteSsid(ssid)
	}
	return ssids
}

func (this *Site) UpdateSSID(ssid *site.SSID) error {
	//TODO implement me
	panic("implement me")
}

func (this *Site) RemoveSSID(name string) error {
	var ix int
	for ix = 0; ix < len(this.ssids) && this.ssids[ix].Name != name; ix++ {
	}
	if ix == len(this.ssids) {
		return errors.New("ssid not found")
	}
	for _, ap := range this.accessPoints {
		if err := ap.RemoveSSID(ssidToSiteSsid(this.ssids[ix])); err != nil {
			// TODO restore SSID on succeeded access points
			return err
		}
	}
	for i := ix + 1; i < len(this.ssids); i++ {
		this.ssids[i-1] = this.ssids[i]
	}
	this.ssids = this.ssids[:len(this.ssids)-1]
	return this.save()
}

func (this *Site) AddStation(ssidName string, station *site.Station) error {
	ix := slices.IndexFunc(this.ssids, func(ssid *SSID) bool {
		return ssidName == ssid.Name
	})
	if ix < 0 {
		return errors.New("SSID \"" + ssidName + "\" not found")
	}
	ssid := this.ssids[ix]
	mac := strings.ToLower(station.Mac)
	stix := slices.IndexFunc(ssid.Stations, func(st *Station) bool {
		return st.Mac == mac
	})
	if stix != -1 {
		ssid.Stations[stix].Name = station.Name
		ssid.Stations[stix].Comment = station.Comment
	} else {
		this.ssids[ix].Stations = append(ssid.Stations, siteStationToStation(station))
	}
	return this.save()
}

func (this *Site) GetStations(ssidName string) ([]*site.Station, error) {
	ix := slices.IndexFunc(this.ssids, func(ssid *SSID) bool {
		return ssidName == ssid.Name
	})
	if ix < 0 {
		return nil, errors.New("SSID \"" + ssidName + "\" not found")
	}
	ssid := this.ssids[ix]
	stations := make([]*site.Station, len(ssid.Stations))
	for i, st := range ssid.Stations {
		stations[i] = stationToSiteStation(st)
	}
	return stations, nil
}

func (this *Site) RemoveStation(ssidName, macAddress string) error {
	ix := slices.IndexFunc(this.ssids, func(ssid *SSID) bool {
		return ssidName == ssid.Name
	})
	if ix < 0 {
		return errors.New("SSID \"" + ssidName + "\" not found")
	}
	ssid := this.ssids[ix]
	mac := strings.ToLower(macAddress)
	stix := slices.IndexFunc(ssid.Stations, func(st *Station) bool {
		return st.Mac == mac
	})
	if stix == -1 {
		return errors.New("Station \"" + mac + "\" not found in " + ssidName + " stations list")
	} else {
		// TODO: remove station at index stix from stations list
	}
	return this.save()
}

func (this *Site) AddDeviceType(device *site.AccessPointDevice) error {
	if this.devices[device.Name] != nil {
		return errors.New("device " + device.Name + " already exists")
	}
	this.devices[device.Name] = siteDeviceToDevice(device)
	return nil
}

func (this *Site) RemoveDeviceType(deviceType string) error {
	for _, ap := range this.accessPoints {
		if ap.Model == deviceType {
			return errors.New("Access points with device type " + deviceType + " exists, device type can't be deleted")
		}
	}
	delete(this.devices, deviceType)
	return nil
}

func (this *Site) GetDeviceTypes() []*site.AccessPointDevice {
	devices := make([]*site.AccessPointDevice, len(this.devices), len(this.devices))
	ix := 0
	for _, dt := range this.devices {
		devices[ix] = deviceToSiteDevice(dt)
		ix++
	}
	return devices
}

func (this *Site) Export(dest io.Writer) error {
	model := this.export()
	content, err := yaml.Marshal(model)
	if err != nil {
		return err
	}
	_, err = dest.Write(content)
	return err
}

func (this *Site) load() error {
	model := new(SiteModel)
	if err := util.ReadObject(this.path, model); err != nil {
		return err
	}
	return this.init(model)
}

func (this *Site) init(model *SiteModel) error {
	this.plugin = pluginName
	this.sshKey = model.SshKey
	this.sshPublicKey = model.SshPublicKey
	this.password = model.Password

	this.ssids = make([]*SSID, len(model.Ssids))
	this.devices = make(map[string]*AccessPointDevice)
	for _, device := range model.Devices {
		if device != nil {
			this.devices[device.Name] = device
		}
	}
	this.accessPoints = make(map[string]*AccessPoint)
	for _, ap := range model.AccessPoints {
		accessPoint, err := NewAccessPoint(ap, this)
		if err != nil {
			return errors.New("AccessPoint " + ap.Name + " has unknown device type " + ap.Model)
		}
		this.accessPoints[ap.Name] = accessPoint
	}
	return nil
}

func (this *Site) export() *SiteModel {
	model := new(SiteModel)
	model.SshKey = this.sshKey
	model.SshPublicKey = this.sshPublicKey
	model.Password = this.password
	model.Devices = make([]*AccessPointDevice, len(this.devices))
	j := 0
	for _, device := range this.devices {
		model.Devices[j] = device
		j++
	}
	model.Ssids = make([]*SSID, len(this.ssids))
	for i, ssid := range this.ssids {
		model.Ssids[i] = new(SSID)
		*model.Ssids[i] = *ssid
	}
	model.AccessPoints = make([]*AccessPointModel, len(this.accessPoints))
	j = 0
	for _, ap := range this.accessPoints {
		model.AccessPoints[j] = ap.export()
		j++
	}
	return model
}

func (this *Site) save() error {
	return util.WriteObject(this.path, this.export())
}
