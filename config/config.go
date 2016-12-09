package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gcfg.v1"
)

var (
	// DefaultConfigFile is the default place json2envdir will look if
	// no config is passed in via a cli flag.
	DefaultConfigFile = "/etc/json2envdir.conf"

	defaultFilePerms os.FileMode = 0640
	defaultPathPerms os.FileMode = 0750
)

// Config holds parsed config.
type Config struct {
	Sections map[string]*EnvDirSection `gcfg:"envdir"`
}

// EnvDirSection represents a configured envdir diretive in the config.
type EnvDirSection struct {
	Path            []string
	PathPermsString string `gcfg:"path-perms"`
	FilePermsString string `gcfg:"file-perms"`

	PathPerms os.FileMode
	FilePerms os.FileMode
}

// LoadConfig loads the specified config file. If an error is encountered the program will panic.
func LoadConfig(configFile string) Config {
	if len(configFile) < 1 {
		configFile = DefaultConfigFile
	}

	cfg := Config{}
	err := gcfg.ReadFileInto(&cfg, configFile)
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetEnv returns the requested EnvDirSection. If non exist it will return an error.
func (cfg Config) GetEnv(env string) (*EnvDirSection, error) {
	if conf, ok := cfg.Sections[env]; ok {
		conf.parsePerms()
		return conf, nil
	}

	return nil, fmt.Errorf("envdir '%s' not found in config", env)
}

func parsePerms(perms string, def os.FileMode) os.FileMode {
	if len(perms) > 0 {
		p, err := strconv.ParseUint(perms, 8, 32)
		if err != nil {
			panic(err)
		}
		return os.FileMode(p)
	}
	return def
}

func (e *EnvDirSection) parsePerms() {
	e.PathPerms = parsePerms(e.PathPermsString, defaultPathPerms)
	e.FilePerms = parsePerms(e.FilePermsString, defaultFilePerms)
}
