package command

import (
	"fmt"
	"strings"
	"wnetctl/config"
	"wnetctl/site"
)

func GetApCommand(argv []string) Command {
	var cmd Command
	switch argv[0] {
	case "add":
		cmd = new(apAdd)
	case "tune":
		cmd = new(apTune)
	case "replace":
		cmd = new(apReplace)
	case "remove":
		cmd = new(apRemove)
	default:
		cmd = apHelp(true)
	}
	cmd.Init()
	if cmd.ParseArgs(argv[1:]) != nil {
		return apHelp(true)
	}
	return cmd
}

type apHelp bool

func (this apHelp) Init() {
}

func (this apHelp) HelpRequested() bool {
	return true
}

func (this apHelp) HelpMessage() string {
	messages := []string{
		"Access Point commands:",
		"add <apName> -t apType -a apIp",
		"tune <apName> [-2c channel] [-2p power] [-5c channel] [-5p power]",
		"replace <apName> -t apType -i apIp",
		"remove <apName>"}
	return strings.Join(messages, "\n  ")
}

func (this apHelp) ParseArgs(argv []string) error {
	return nil
}

func (this apHelp) Execute() error {
	fmt.Println(this.HelpMessage())
	return nil
}

type apCommand struct {
	GenericCommand
	name string
}

type apAdd struct {
	apCommand
	model *site.AccessPointRequest
}

func (this *apAdd) Init() {
	this.GenericCommand.Init()
	this.usageMessage = "Usage: wnetctl ap add <apName> -t apType -a apIp"
	this.model = new(site.AccessPointRequest)
	this.flags.StringVar(&this.model.Ip, "a", "", "access point IP address")
	this.flags.StringVar(&this.model.Ip, "addr", "", "access point IP address")
	this.flags.StringVar(&this.model.Model, "t", "", "access point device type")
	this.flags.StringVar(&this.model.Model, "type", "", "access point device type")
}

func (this *apAdd) ParseArgs(argv []string) error {
	err := this.flags.Parse(argv)
	if err != nil {
		return err
	}
	if this.flags.NArg() < 1 {
		this.helpRequested = true
	} else {
		this.model.Name = this.flags.Arg(0)
	}
	return nil
}

func (this *apAdd) Execute() error {
	siteManager, err := config.GetCurrentSiteManager(getSiteManager)
	if err == nil {
		_, err = siteManager.AddAccessPoint(this.model)
	}
	return err
}

type apTune struct {
	apCommand
}

func (this *apTune) HelpMessage() string {
	//TODO implement me
	panic("implement me")
}

func (this *apTune) ParseArgs(argv []string) error {
	//TODO implement me
	panic("implement me")
}

func (this *apTune) Execute() error {
	//TODO implement me
	panic("implement me")
}

type apRemove struct {
	GenericCommand
	names []string
}

func (this *apRemove) Init() {
	this.GenericCommand.Init()
	this.usageMessage = "Usage: wnetctl ap remove <apName ...>"
}

func (this *apRemove) ParseArgs(argv []string) error {
	err := this.flags.Parse(argv)
	if err != nil {
		this.helpRequested = true
		return err
	}
	if this.flags.NArg() < 1 {
		this.helpRequested = true
		return nil
	}
	this.names = make([]string, this.flags.NArg())
	copy(this.names, this.flags.Args()[0:this.flags.NArg()])
	return nil
}

func (this *apRemove) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	siteManager, err := config.GetCurrentSiteManager(getSiteManager)
	if err != nil {
		return err
	}
	for _, name := range this.names {
		err = siteManager.RemoveAccessPoint(name)
		if err != nil {
			return err
		}
	}
	return nil
}

type apReplace struct {
	apCommand
}

func (this *apReplace) HelpMessage() string {
	//TODO implement me
	panic("implement me")
}

func (this *apReplace) ParseArgs(argv []string) error {
	//TODO implement me
	panic("implement me")
}

func (this *apReplace) Execute() error {
	//TODO implement me
	panic("implement me")
}
