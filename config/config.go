package config

import (
	"errors"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"wnetctl/site"
	"wnetctl/util"
)

const configDir = "wnetctl"
const configFile = "config.yml"
const sitesDir = "sites"
const defaultType = "openwrt"

type SiteInfo struct {
	Name        string
	Description string
	Filepath    string
	Type        string
}

type SitesConfig interface {
	List() []*SiteInfo
	Select(name string) (*SiteInfo, error)
	Current() *SiteInfo
	Add(name, description string) (*SiteInfo, error)
	Remove(name string) (*SiteInfo, error)
}

type sitesConfig struct {
	SitesInfo    []*SiteInfo "yaml:\"sites\""
	SelectedSite string      "yaml:\"selected\""
	sites        map[string]*SiteInfo
}

var defaultSitesConfig *sitesConfig

func (this *sitesConfig) List() []*SiteInfo {
	sitesInfo := make([]*SiteInfo, len(this.SitesInfo))
	for i, info := range this.SitesInfo {
		sitesInfo[i] = &SiteInfo{Name: info.Name, Description: info.Description, Filepath: info.Filepath}
	}
	return sitesInfo
}

func (this *sitesConfig) Select(name string) (*SiteInfo, error) {
	siteInfo := this.sites[name]
	if siteInfo == nil {
		return nil, errors.New("Site not found: " + name)
	}
	this.SelectedSite = name
	if err := saveSitesConfig(this); err != nil {
		return nil, err
	}
	return &SiteInfo{Name: siteInfo.Name, Description: siteInfo.Description, Filepath: siteInfo.Filepath}, nil
}

func (this *sitesConfig) Current() *SiteInfo {
	if len(this.sites) == 0 {
		return nil
	}
	siteInfo := this.sites[this.SelectedSite]
	return &SiteInfo{Name: siteInfo.Name, Description: siteInfo.Description, Filepath: siteInfo.Filepath}
}

func (this *sitesConfig) Add(name, description string) (*SiteInfo, error) {
	if this.sites[name] != nil {
		return nil, errors.New("Site already exists: " + name)
	}
	configFilepath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configFilepath = filepath.Join(configFilepath, configDir, sitesDir)
	if err = verifyDirectory(configFilepath); err != nil {
		return nil, err
	}
	siteFileName := filepath.Join(configFilepath, uuid.NewString()+".yml")
	siteInfo := &SiteInfo{Name: name, Description: description, Filepath: siteFileName}
	if len(this.SitesInfo) == 0 {
		this.SelectedSite = name
	}
	this.SitesInfo = append(this.SitesInfo, siteInfo)
	this.sites[name] = siteInfo
	if err := saveSitesConfig(this); err != nil {
		return nil, err
	}
	return siteInfo, nil
}

func (this *sitesConfig) Remove(name string) (*SiteInfo, error) {
	var newSelected *SiteInfo
	if len(this.SitesInfo) == 1 {
		newSelected = nil
	} else {
		if name == this.SelectedSite {

			//newSelected =
		}
	}
	return newSelected, nil
}

func loadSitesConfig(scfg *sitesConfig) error {
	configFilepath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configFilepath = filepath.Join(configFilepath, configDir)
	if err = verifyDirectory(configFilepath); err != nil {
		return err
	}
	configFilepath = filepath.Join(configFilepath, configFile)
	fileInfo, err := os.Stat(configFilepath)
	if err != nil { // FIXME ensure err is about file does not exist
		scfg.sites = make(map[string]*SiteInfo)
		return nil
	} else {
		if fileInfo.IsDir() {
			return errors.New("Sites config path is a directory but should be a file: " + configFilepath)
		}
	}
	if err = util.ReadObject(configFilepath, scfg); err != nil {
		return err
	}
	for _, siteInfo := range scfg.SitesInfo {
		scfg.sites[siteInfo.Name] = siteInfo
		if siteInfo.Type == "" {
			siteInfo.Type = defaultType
		}
	}
	return nil
}

func saveSitesConfig(scfg *sitesConfig) error {
	configFilepath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configFilepath = filepath.Join(configFilepath, configDir)
	if err = verifyDirectory(configFilepath); err != nil {
		return err
	}
	configFilepath = filepath.Join(configFilepath, configFile)
	return util.WriteObject(configFilepath, scfg)
}

func GetSitesConfig() (SitesConfig, error) {
	if defaultSitesConfig == nil {
		config := new(sitesConfig)
		config.sites = make(map[string]*SiteInfo)
		defaultSitesConfig = config
		if err := loadSitesConfig(config); err != nil {
			return nil, err
		}
	}
	return defaultSitesConfig, nil
}

func verifyDirectory(dir string) error {
	fileInfo, err := os.Stat(dir)
	if err == nil {
		if fileInfo.IsDir() {
			return nil
		}
		return errors.New(dir + " exists but is not a directory: " + dir)
	} else {
		return os.MkdirAll(dir, 0700)
	}
}

func GetCurrentSiteManager(factory func(string, string, string) (site.SiteManager, error)) (site.SiteManager, error) {
	sitesConfig, err := GetSitesConfig()
	if err != nil {
		return nil, err
	}
	siteInfo := sitesConfig.Current()
	if siteInfo == nil {
		return nil, errors.New("There are no sites yet, nothing to do")
	}
	siteType := siteInfo.Type
	if siteType == "" {
		siteType = defaultType
	}
	return factory(siteType, siteInfo.Name, siteInfo.Filepath)
}
