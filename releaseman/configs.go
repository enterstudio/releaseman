package releaseman

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"gopkg.in/yaml.v2"
)

//=======================================
// Consts
//=======================================

const (
	// DefaultConfigPth ...
	DefaultConfigPth = "./release_config.yml"
)

var (
	// IsCIMode ...
	IsCIMode = false
)

//=======================================
// Models
//=======================================

// Release ...
type Release struct {
	StartFromBranch string   `yaml:"start_from"`
	ReleaseOnBranch string   `yaml:"release_on"`
	Version         string   `yaml:"version,omitempty"`
	Changes         []string `yaml:"changes,omitempty"`
}

// Changelog ...
type Changelog struct {
	Path            string `yaml:"path"`
	ContentTemplate string `yaml:"content_template"`
	HeaderTemplate  string `yaml:"header_template"`
	FooterTemplate  string `yaml:"footer_template"`
}

// Config ...
type Config struct {
	Release   Release   `yaml:"release,omitempty"`
	Changelog Changelog `yaml:"changelog,omitempty"`
}

// NewConfigFromFile ...
func NewConfigFromFile(pth string) (Config, error) {
	bytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return Config{}, err
	}
	return NewConfigFromBytes(bytes)
}

// NewConfigFromBytes ...
func NewConfigFromBytes(bytes []byte) (Config, error) {
	type FileConfig struct {
		Release   *Release   `yaml:"release,omitempty"`
		Changelog *Changelog `yaml:"changelog,omitempty"`
	}

	fileConfig := FileConfig{}
	if err := yaml.Unmarshal(bytes, &fileConfig); err != nil {
		return Config{}, err
	}

	config := Config{}
	if fileConfig.Release != nil {
		config.Release = *fileConfig.Release
	}
	if fileConfig.Changelog != nil {
		config.Changelog = *fileConfig.Changelog
	}

	return config, nil
}

// WriteConfigToFile ...
func WriteConfigToFile(config Config, pth string) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return fileutil.WriteBytesToFile(pth, bytes)
}

// PrintMode ...
type PrintMode uint8

const (
	// FullMode ...
	FullMode PrintMode = iota
	// ChangelogMode ...
	ChangelogMode
	// ReleaseMode ...
	ReleaseMode
)

// Print ...
func (config Config) Print(mode PrintMode) {
	fmt.Println()
	log.Infof("Your configuration:")

	if mode == ChangelogMode || mode == ReleaseMode || mode == FullMode {
		log.Infof(" * Start from branch: %s", config.Release.StartFromBranch)
	}
	if mode == ReleaseMode || mode == FullMode {
		log.Infof(" * Release on branch: %s", config.Release.ReleaseOnBranch)
	}
	if config.Release.Version != "" && (mode == ChangelogMode || mode == ReleaseMode || mode == FullMode) {
		log.Infof(" * Release version: %s", config.Release.Version)
	}
	if mode == ChangelogMode || mode == FullMode {
		log.Infof(" * Changelog path: %s", config.Changelog.Path)
	}

	fmt.Println()
}

func getReleasemanDirPath() (string, error) {
	pth := path.Join(pathutil.UserHomeDir(), ".releaseman")
	return pth, EnsureDirExists(pth)
}

// GetPreparedReleaseConfigPath ...
func GetPreparedReleaseConfigPath() string {
	return path.Join(getReleasemanDirPath(), "prepared_config.yml")
}

// EnsureDirExists ...
func EnsureDirExists(pth string) error {
	confDirPth := getReleasemanDirPath()
	return pathutil.EnsureDirExist(pth)
}
