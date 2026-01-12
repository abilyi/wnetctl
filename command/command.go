package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"wnetctl/openwrt"
	"wnetctl/site"
)

type Object interface {
	Init()
}

type Command interface {
	Object
	ParseArgs(argv []string) error
	Execute() error
	HelpRequested() bool
	HelpMessage() string
}

type GenericCommand struct {
	helpRequested bool
	usageMessage  string
	//helpMessage   string
	flags *flag.FlagSet
}

type Help bool

func (this Help) HelpRequested() bool {
	return true
}

func (this Help) HelpMessage() string {
	help := []string{"Usage: wnetctl <object> <command> <options>",
		"where <object> is one of: site, device, ap, ssid, station",
		"commands are object specific, although \"help\" command supported for each object explaining available commands",
		"also each command has any of -h, -help --help options with details about options and parameters"}
	return strings.Join(help, "\n  ")
}

func (this Help) ParseArgs(argv []string) error {
	return nil
}

func (this Help) Execute() error {
	fmt.Println(this.HelpMessage())
	return nil
}

func (this Help) Init() {
}

func (this *GenericCommand) getStdFlags(name string) *flag.FlagSet {
	this.flags = flag.NewFlagSet(name, flag.ContinueOnError)
	this.flags.BoolVar(&this.helpRequested, "h", false, "Display help message")
	this.flags.BoolVar(&this.helpRequested, "help", false, "Display help message")
	return this.flags
}

func (this *GenericCommand) HelpRequested() bool {
	return this.helpRequested
}

func (this *GenericCommand) HelpMessage() string {
	sb := new(strings.Builder)
	sb.WriteString(this.usageMessage)
	sb.WriteByte('\n')
	defaultOutput := this.flags.Output()
	this.flags.SetOutput(sb)
	this.flags.PrintDefaults()
	this.flags.SetOutput(defaultOutput)
	return sb.String()
}

func (this *GenericCommand) Init() {
	this.flags = flag.NewFlagSet("commandFlags", flag.ContinueOnError)
	this.flags.BoolVar(&this.helpRequested, "h", false, "Display help message")
	this.flags.BoolVar(&this.helpRequested, "help", false, "Display help message")
}

func getSiteManager(siteType, name, filepath string) (site.SiteManager, error) {
	switch siteType {
	case "openwrt":
		return openwrt.NewSiteManager(name, filepath)
	default:
		return nil, errors.New("Unsupported site type " + siteType)
	}
}

func createSiteManager(siteType, name, filepath string, model *site.SiteRequest) (site.SiteManager, error) {
	switch siteType {
	case "openwrt":
		return openwrt.CreateSiteManager(name, filepath, model)
	default:
		return nil, errors.New("Unsupported site type " + siteType)
	}
}
