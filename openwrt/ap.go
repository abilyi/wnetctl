package openwrt

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"path/filepath"
	"strings"
	"text/template"
	"wnetctl/site"
	"wnetctl/sshclient"
)

const TEMPLATES = "templates/openwrt/"

type AccessPointDevice struct {
	Name               string
	Model              string
	Wlan2              *DeviceWirelessAdapter
	Wlan5              *DeviceWirelessAdapter
	BridgedWiredDevice string
	Architecture       string
	Cpu                string
}

type DeviceWirelessAdapter struct {
	Interface string
	Device    string
	Driver    string
}

type WirelessAdapter struct {
	Device  *DeviceWirelessAdapter
	Mac     string
	Channel int
	Power   int
}

type AccessPoint struct {
	name  string
	Model string
	Mac   string
	Ip    string
	Wlan2 *WirelessAdapter
	Wlan5 *WirelessAdapter
	site  *Site
}

type WirelessAdapterModel struct {
	Mac     string
	Channel int
	Power   int
}

type AccessPointModel struct {
	Name  string
	Model string
	Mac   string
	Ip    string
	WLan2 *WirelessAdapterModel
	WLan5 *WirelessAdapterModel
}

type SSIDModel struct {
}

const accessPointAdmin = "root"
const defaultChannel2G = 6
const defaultChannel5G = 40

func NewWirelessAdapter() *WirelessAdapter {
	return new(WirelessAdapter)
}

func CreateAccessPoint(request *site.AccessPointRequest, site *Site) (*AccessPoint, error) {
	device, exists := site.devices[request.Model]
	if !exists {
		return nil, errors.New("Access point type " + request.Model + " does not exist")
	}
	ap := AccessPoint{site: site, name: request.Name, Model: device.Name, Mac: request.Mac, Ip: request.Ip}
	ap.Wlan2.Device = device.Wlan2
	ap.Wlan5.Device = device.Wlan5
	ap.Wlan5.Channel = defaultChannel5G
	// TODO : discover MAC addresses for both wired and wireless adapters
	return &ap, nil
}

func NewAccessPoint(model *AccessPointModel, site *Site) (*AccessPoint, error) {
	device, exists := site.devices[model.Model]
	if !exists {
		return nil, errors.New("Access point type " + model.Model + " does not exist")
	}
	ap := AccessPoint{site: site, name: model.Name, Model: device.Name, Mac: model.Mac, Ip: model.Ip}
	if model.WLan2 != nil {
		ap.Wlan2 = new(WirelessAdapter)
		modelToWirelessAdapter(ap.Wlan2, model.WLan2, device.Wlan2)
	}
	if model.WLan5 != nil {
		ap.Wlan5 = new(WirelessAdapter)
		modelToWirelessAdapter(ap.Wlan5, model.WLan5, device.Wlan5)
	}
	return &ap, nil
}

func (this *AccessPoint) Name() string {
	return this.name
}

func (this *AccessPoint) Configure() error {
	sshClient := sshclient.NewSshClient(this.Ip, accessPointAdmin, "", "")
	if err := installSshPublicKey(sshClient, this.site.sshPublicKey); err != nil {
		return err
	}
	sshClient.SetKey(this.site.sshKey)
	if err := sshClient.Connect(); err != nil {
		return err
	}
	if err := sshClient.ExecuteInteractive(sshclient.NewPasswd(accessPointAdmin, "", this.site.password)); err != nil {
		return err
	}
	// TODO disable password auth for SSH, set TZ, enable NTP; render initial template
	// TODO install usteer (optional?)
	// TODO render SSID template for each SSID defined
	/*
		scriptBuilder := new(strings.Builder)
		for _, ssid := range this.site.ssids {
			ssidModel := buildSSIDModel(this.site, this.name)
			if err := renderCommands(scriptBuilder, "add_ssid", ssidModel); err != nil {
				return err
			}
		}
		commands := strings.Split(scriptBuilder.String(), "\n")
		for _, command := range commands {
			if err := sshClient.Execute(command); err != nil {
				return err
			}
		}
	*/
	return nil
}

func installSshPublicKey(sshClient sshclient.SshClient, pubkeyPath string) error {
	err := sshClient.Connect()
	if err != nil {
		return err
	}
	err = sshClient.Execute("/bin/test -d /root/.sshclient")
	if err != nil {
		_, exitError := err.(*ssh.ExitError)
		if exitError {
			if err = sshClient.Execute("/bin/mkdir /root/.sshclient"); err != nil {
				return err
			}
			if err = sshClient.Execute("/bin/chmod 0700 /root/.sshclient"); err != nil {
				return err
			}
		}
	}
	if err = sshClient.ExecuteInteractive(sshclient.NewInstallSshKey(pubkeyPath)); err != nil {
		return err
	}
	if err = sshClient.Execute("/bin/chmod 0600 /root/.sshclient/authorized_keys"); err != nil {
		return err
	}
	return sshClient.Close()
}

func (this *AccessPoint) AddNeighbour(neighbour site.AccessPoint) error {
	// var commands, rollbackCmds []string

	// commands = append(commands, "uci add-list ... "+buildNeighborRoamingValue(neighbour))
	// TODO create model object

	return nil
	// return sshClient.Execute(this.IPAddress, commands, rollbackCmds)
}

func (this *AccessPoint) RemoveNeighbour(neighbour site.AccessPoint) error {
	//TODO implement me
	panic("implement me")
}

func (this *AccessPoint) AddSSID(ssid *site.SSID) error {
	// var commands, rollbackCmds []string
	// TODO create model object
	//TODO add commands creating vlan if any

	// return this.sshClient.Execute(this.IPAddress, commands, rollbackCmds)
	panic("implement me")
}

func (this *AccessPoint) RemoveSSID(ssid *site.SSID) error {
	//TODO implement me
	panic("implement me")
}

func (this *AccessPoint) AddStation(mac string) error {
	return nil
}

func (this *AccessPoint) RemoveStation(mac string) error {
	return nil
}

func (this *AccessPoint) ToResponse() *site.AccessPointResponse {
	model := new(site.AccessPointResponse)
	model.Name = this.name
	model.Model = this.Model
	model.Mac = this.Mac
	model.Ip = this.Ip
	return model
}

func (this *AccessPoint) export() *AccessPointModel {
	model := new(AccessPointModel)
	model.Name = this.name
	model.Model = this.Model
	model.Mac = this.Mac
	model.Ip = this.Ip
	return model
}

func buildSSIDModel(site *Site, targetAp string) *SSIDModel {
	model := new(SSIDModel)
	// TODO
	return model
}

func buildNeighborRoamingValue(point *AccessPoint) string {
	// TODO complete parameter generation
	return strings.ReplaceAll(point.Mac, ":", "_")
}

func renderCommands(out io.Writer, scriptTemplate string, data interface{}) error {
	gotmpl := template.Must(template.New(scriptTemplate).Parse(filepath.Clean(TEMPLATES + scriptTemplate)))
	return gotmpl.Execute(out, data)
}
