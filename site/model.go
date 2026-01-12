package site

import (
	"fmt"
	"strings"
)

/*
type SiteListEntry struct {
	Name        string
	Filename    string
	Description string
}

type SitesList struct {
	sites map[string]SiteListEntry
	selected
}
*/

type DeviceWirelessAdapter struct {
	Interface string
	Device    string
	Driver    string
}

type WirelessAdapterModel struct {
	DeviceWirelessAdapter
	Channel int
	Power   int
}

type AccessPointRequest struct {
	Name  string
	Model string
	Mac   string
	Ip    string
}

type AccessPointResponse struct {
	AccessPointRequest
	WLan2 *WirelessAdapterModel
	WLan5 *WirelessAdapterModel
}

type AccessPointDevice struct {
	Name               string
	Model              string
	WLan2              *DeviceWirelessAdapter
	WLan5              *DeviceWirelessAdapter
	BridgedWiredDevice string
	Architecture       string
	Cpu                string
}

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

type SiteRequest struct {
	SshKey       string `yaml:"sshKey"`
	SshPublicKey string `yaml:"sshPublicKey"`
	Password     string
	Country      string
	SsidSuffix2  string `yaml:"ssidSuffix2"`
	SsidSuffix5  string `yaml:"ssidSuffix5"`
}

type SiteResponse struct {
	SshKey       string `yaml:"sshKey"`
	SshPublicKey string `yaml:"sshPublicKey"`
	Password     string
	Country      string
	SsidSuffix2  string                `yaml:"ssidSuffix2"`
	SsidSuffix5  string                `yaml:"ssidSuffix5"`
	AccessPoints []*AccessPointRequest `yaml:"accessPoints"`
	Ssid         []*SSID
	Devices      []*AccessPointDevice
}

func (this *AccessPointRequest) String() string {
	mac := "unknown"
	if this.Mac != "" {
		mac = this.Mac
	}
	//return fmt.Sprintf("AP %s (%s) IP %s MAC %s. 2.4GHz radio on channel %d at %d dBm, 5GHz radio on channel %d at %d dBm",
	//	this.Name, this.Model, this.Ip, mac, this.WLan2.Channel, this.WLan2.Power, this.WLan5.Channel, this.WLan5.Power)
	return fmt.Sprintf("AP %s (%s) IP %s, MAC %s",
		this.Name, this.Model, this.Ip, mac)
}

func (this DeviceWirelessAdapter) String() string {
	return fmt.Sprintf("%s (%s) driver %s", this.Device, this.Interface, this.Driver)
}

func (this *AccessPointDevice) String() string {
	info := []string{}
	info = append(info, fmt.Sprintf("AP model: %s, short name %s, architecture %s, CPU %s.", this.Model, this.Name, this.Architecture, this.Cpu))
	if this.WLan2 != nil {
		info = append(info, fmt.Sprintf("2.4GHz WiFi: %s", this.WLan2))
	} else {
		info = append(info, "No 2.4GHz WiFi")
	}
	if this.WLan5 != nil {
		info = append(info, fmt.Sprintf("5GHz WiFi: %s", this.WLan5))
	} else {
		info = append(info, "No 5GHz WiFi")
	}
	info = append(info, fmt.Sprintf("Wired interface %s", this.BridgedWiredDevice))
	return strings.Join(info, "\n  ")
}

func NewSSID() *SSID {
	return new(SSID)
}

func (this *SSID) String() string {
	if this.Vlan > 0 {
		return fmt.Sprintf("%s on vlan %d, auth %s", this.Name, this.Vlan, this.Auth)
	} else {
		return fmt.Sprintf("%s on default vlan, auth %s", this.Name, this.Auth)
	}
}

func (this *SiteResponse) String() string {
	info := []string{}
	info = append(info, fmt.Sprintf("Site configuration:\nSSH key: %s (public %s)", this.SshKey, this.SshPublicKey))
	info = append(info, fmt.Sprintf("2.4GHz wlan networks suffix: \"%s\"; 5GHz wlan networks suffix: \"%s\"", this.SsidSuffix2, this.SsidSuffix5))
	info = append(info, "* Access points:")
	for _, ap := range this.AccessPoints {
		info = append(info, ap.String())
	}
	info = append(info, "* SSID (wireless networks)")
	for _, ss := range this.Ssid {
		info = append(info, ss.String())
	}
	info = append(info, "* Known device types")
	for _, dev := range this.Devices {
		info = append(info, dev.String())
	}
	return strings.Join(info, "\n")
}
