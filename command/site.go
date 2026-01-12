package command

import (
	"fmt"
	"io"
	"os"
	"strings"
	"wnetctl/config"
	"wnetctl/site"
)

func GetSiteCommand(argv []string) Command {
	var cmd Command
	switch argv[0] {
	case "list":
		cmd = new(siteList)
	case "init":
		cmd = new(siteInit)
	case "select":
		cmd = new(siteSelect)
	case "export":
		cmd = new(siteExport)
	case "import":
		cmd = new(siteImport)
	default:
		cmd = siteHelp(true)
	}
	if cmd.ParseArgs(argv[1:]) != nil {
		return siteHelp(true)
	}
	return cmd
}

type SiteCommand struct {
	GenericCommand
	name string
}

func (this *SiteCommand) Init() {
	this.GenericCommand.Init()
}

type siteList struct {
	GenericCommand
}

func (this *siteList) Init() {
	this.GenericCommand.Init()
	this.usageMessage = "Usage: wnetctl site list\nNo options, no arguments. Current site is marked with asterisk"
}

func (this *siteList) ParseArgs(argv []string) error {
	this.helpRequested = len(argv) > 0
	return nil
}

func (this *siteList) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	cfg, err := config.GetSitesConfig()
	if err != nil {
		return err
	}
	sites := cfg.List()
	current := cfg.Current()
	for _, s := range sites {
		if current.Name == s.Name {
			fmt.Printf("* %s\t%s\n", s.Name, s.Description)
		} else {
			fmt.Printf("  %s\t%s\n", s.Name, s.Description)
		}
	}
	return nil
}

type siteSelect struct {
	SiteCommand
}

func (this *siteSelect) Init() {
	this.SiteCommand.Init()
	this.usageMessage = "Usage: wnetctl site select <name>"
}

func (this *siteSelect) ParseArgs(argv []string) error {
	if len(argv) != 1 {
		this.helpRequested = true
	}
	flags := this.getStdFlags("siteExport")
	if err := flags.Parse(argv); err != nil {
		return err
	}
	this.helpRequested = this.helpRequested || flags.NArg() != 1
	this.name = flags.Arg(0)
	return nil
}

func (this *siteSelect) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	cfg, err := config.GetSitesConfig()
	if err != nil {
		return err
	}
	_, err = cfg.Select(this.name)
	return err
}

type siteInit struct {
	SiteCommand
	siteModel   *site.SiteRequest
	description string
}

func (this *siteInit) Init() {
	this.SiteCommand.Init()
	this.siteModel = new(site.SiteRequest)
	this.flags.StringVar(&this.siteModel.SshKey, "sk", "", "path to SSH private key")
	this.flags.StringVar(&this.siteModel.SshPublicKey, "sp", "", "path to SSH public key")
	this.flags.StringVar(&this.siteModel.Password, "p", "", "root (or admin) password for access points")
	this.flags.StringVar(&this.siteModel.SsidSuffix2, "s2", "", "suffix appended to any SSID in 2.4 GHz band")
	this.flags.StringVar(&this.siteModel.SsidSuffix5, "s5", "", "suffix appended to any SSID in 5 GHz band")

	this.usageMessage = "Usage: wnetctl site init site_name <options>"
}

func (this *siteInit) ParseArgs(argv []string) error {
	if len(argv) == 0 {
		this.helpRequested = true
		return nil
	}
	err := this.flags.Parse(argv[1:])
	if err != nil {
		return err
	}
	this.name = argv[0]
	if this.helpRequested {
		return nil
	}
	// TODO validate required fields set and keys exists
	return nil
}

func (this *siteInit) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	cfg, err := config.GetSitesConfig()
	if err != nil {
		return err
	}
	siteInfo, err := cfg.Add(this.name, this.description)
	if err != nil {
		return err
	}
	/*siteManager*/ _, err = getSiteManager(siteInfo.Type, siteInfo.Name, siteInfo.Filepath)
	if err != nil {
		return err
	}
	//siteManager.
	/*	finalModel := site.NewSiteModel()
		if err := siteManager.Export(finalModel); err != nil {
			return err
		}
		cfg, err := config.GetSitesConfig()
		if err != nil {
			return err
		}
		return util.WriteObject(siteInfo.Filepath, finalModel)
	*/
	return nil
}

type siteExport struct {
	SiteCommand
	filename string
}

func (this *siteExport) Init() {
	this.SiteCommand.Init()
	this.usageMessage = "Usage: wnetctl site export [filename]"
}

func (this *siteExport) ParseArgs(argv []string) error {
	if err := this.flags.Parse(argv); err != nil {
		this.helpRequested = true
	}
	if this.flags.NArg() != 1 {
		this.helpRequested = true
	} else {
		this.filename = this.flags.Args()[0]
	}
	return nil
}

func (this *siteExport) Execute() error {
	if this.helpRequested {
		fmt.Println(this.HelpMessage())
		return nil
	}
	cfg, err := config.GetSitesConfig()
	if err != nil {
		return err
	}
	siteInfo := cfg.Current()
	if siteInfo == nil {
		return fmt.Errorf("There is no current site.")
	}
	in, err := os.OpenFile(siteInfo.Filepath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	} else {
		defer in.Close()
	}
	var out io.Writer
	if this.filename == "" {
		out = os.Stdout
	} else {
		file, err := os.OpenFile(this.filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		out = file
	}
	_, err = io.Copy(out, in)
	return err
}

type siteImport struct {
	SiteCommand
	description string
	filename    string
}

func (this *siteImport) Init() {
	this.SiteCommand.Init()
	//TODO implement me
	panic("implement me")
}

func (this *siteImport) ParseArgs(argv []string) error {
	//TODO implement me
	panic("implement me")
}

func (this *siteImport) Execute() error {
	//TODO implement me
	panic("implement me")
}

type siteHelp bool

func (siteHelp) Init() {
}

func (this siteHelp) HelpRequested() bool {
	return true
}

func (this siteHelp) HelpMessage() string {
	help := []string{"Usage: wnetctl site <command> <params...>\nSupported commands are:",
		"list    Show known sites. Has no parameters and options except -h, --help, -help",
		"init    Create an initialized site configuration file. For more details use wnetctl site init -h",
		"select  Selects site so any further commands are applied to it. For more details use wnetctl site select -h",
		"export  Exports site configuration file. For more details use wnetctl site export -h",
		"import  Imports site configuration file. For more details use wnetctl site import -h",
		"help    Show this help text."}
	return strings.Join(help, "\n  ")
}

func (this siteHelp) ParseArgs(argv []string) error {
	return nil
}

func (this siteHelp) Execute() error {
	fmt.Println(this.HelpMessage())
	return nil
}
