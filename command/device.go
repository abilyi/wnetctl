package command

import (
	"fmt"
	"os"
	"strings"
	"wnetctl/config"
	"wnetctl/site"
	"wnetctl/util"
)

func GetDeviceCommand(argv []string) Command {
	var cmd Command
	switch argv[0] {
	case "add":
		cmd = new(deviceAdd)
	case "remove":
		cmd = new(deviceRemove)
	case "list":
		cmd = new(deviceList)
	default:
		cmd = deviceHelp(true)
	}
	cmd.Init()
	if cmd.ParseArgs(argv[1:]) != nil {
		cmd = deviceHelp(true)
	}
	return cmd
}

type deviceHelp bool

func (this deviceHelp) Init() {
}

func (this deviceHelp) ParseArgs(argv []string) error {
	return nil
}

func (this deviceHelp) Execute() error {
	fmt.Println(this.HelpMessage())
	return nil
}

func (this deviceHelp) HelpRequested() bool {
	return true
}

func (this deviceHelp) HelpMessage() string {
	//return string(this)
	help := []string{"Usage: wnetctl device <command> [options]\nAvailable commands are:",
		"add <options>", "remove <name>", "list", "help"}
	msg := strings.Join(help, "\n  ")
	help = []string{msg, "Use wnetctl device <command> -h for details about distinct command."}
	return strings.Join(help, "\n")

}

type deviceCommand struct {
	GenericCommand
	name string
}

type deviceAdd struct {
	deviceCommand
	device     *site.AccessPointDevice
	importPath string
}

func (this *deviceAdd) Init() {
	this.GenericCommand.Init()
	this.device = new(site.AccessPointDevice)
	this.device.WLan2 = new(site.DeviceWirelessAdapter)
	this.device.WLan5 = new(site.DeviceWirelessAdapter)

	this.flags.StringVar(&this.device.Name, "n", "", "Short name of the device")
	this.flags.StringVar(&this.device.Name, "name", "", "Short name of the device")
	this.flags.StringVar(&this.device.Model, "m", "", "Device model")
	this.flags.StringVar(&this.device.Model, "model", "", "Device model")
	this.flags.StringVar(&this.device.Architecture, "a", "", "Device CPU architecture")
	this.flags.StringVar(&this.device.Architecture, "arch", "", "Device CPU architecture")
	this.flags.StringVar(&this.device.Cpu, "c", "", "Device CPU model")
	this.flags.StringVar(&this.device.Cpu, "cpu", "", "Device CPU model")
	this.flags.StringVar(&this.device.BridgedWiredDevice, "e", "", "Device's bridged wired interface")
	this.flags.StringVar(&this.device.BridgedWiredDevice, "eth", "", "Device's bridged wired interface")
	this.flags.StringVar(&this.device.WLan2.Device, "w2dev", "", "2.4GHz WLAN network device (i.e. wlan0)")
	this.flags.StringVar(&this.device.WLan2.Interface, "w2if", "", "2.4GHz WLAN interface (i.e. radio0)")
	this.flags.StringVar(&this.device.WLan2.Driver, "w2drv", "", "2.4GHz WLAN network device driver (i.e. linux module name)")
	this.flags.StringVar(&this.device.WLan2.Device, "w5dev", "", "5GHz WLAN network device (i.e. wlan0)")
	this.flags.StringVar(&this.device.WLan2.Interface, "w5if", "", "5GHz WLAN interface (i.e. radio0)")
	this.flags.StringVar(&this.device.WLan2.Driver, "w5drv", "", "5GHz WLAN network device driver (i.e. linux module name)")
	this.flags.StringVar(&this.importPath, "i", "", "Import from a file containing device information instead of passing values in options")
	this.flags.StringVar(&this.importPath, "import", "", "Path to a file containing device information instead of passing values in options")

	this.usageMessage = "Usage: wnetctl device add <options>"
}

func (this *deviceAdd) ParseArgs(argv []string) error {
	this.device = new(site.AccessPointDevice)

	if err := this.flags.Parse(argv); err != nil {
		this.helpRequested = true
	}
	return nil
}

func (this *deviceAdd) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	if this.importPath != "" {
		if err := util.ReadObject(this.importPath, this.device); err != nil {
			return err
		}
	}
	siteManager, err := config.GetCurrentSiteManager(getSiteManager)
	if err != nil {
		return err
	}
	return siteManager.AddDeviceType(this.device)
}

func (this *deviceAdd) HelpMessage() string {
	msg := new(strings.Builder)
	msg.WriteString("Usage: wnetctl device add <options>\n")
	this.flags.SetOutput(msg)
	this.flags.PrintDefaults()
	this.flags.SetOutput(os.Stderr)
	message := msg.String()
	return strings.Join(strings.Split(message, "\n"), "\n  ")
}

type deviceRemove struct {
	deviceCommand
}

func (this *deviceRemove) Init() {
	this.GenericCommand.Init()
	this.usageMessage = "Usage: wnetctl device remove <name>\n"
}

func (this *deviceRemove) ParseArgs(argv []string) error {
	if err := this.flags.Parse(argv); err != nil {
		this.helpRequested = true
		return nil
	}
	if this.flags.NArg() != 1 {
		this.helpRequested = true
	} else {
		this.name = this.flags.Arg(0)
	}
	return nil
}

func (this *deviceRemove) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	siteManager, err := config.GetCurrentSiteManager(getSiteManager)
	if err != nil {
		return err
	}
	return siteManager.RemoveDeviceType(this.name)
}

type deviceList struct {
	GenericCommand
}

func (this *deviceList) Init() {
	this.GenericCommand.Init()
	this.usageMessage = "Usage: wnetctl device list\n"
}

func (this *deviceList) ParseArgs(argv []string) error {
	if err := this.flags.Parse(argv); err != nil {
		this.helpRequested = true
	}
	if this.flags.NArg() > 0 {
		this.helpRequested = true
	}
	return nil
}

func (this *deviceList) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	siteManager, err := config.GetCurrentSiteManager(getSiteManager)
	if err != nil {
		return err
	}
	devices := siteManager.GetDeviceTypes()
	for _, device := range devices {
		fmt.Println(device.String())
	}
	return nil
}
