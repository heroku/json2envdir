package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gcfg.v1"
)

var (
	DefaultConfigFile             = "/etc/json2envdir.conf"
	DefaultFilePerms  os.FileMode = 0640
	DefaultPathPerms  os.FileMode = 0750
)

type Config struct {
	Sections map[string]*EnvDirSection `gcfg:"envdir"`
}

type EnvDirSection struct {
	Path            []string
	PathPermsString string `gcfg:"path-perms"`
	FilePermsString string `gcfg:"file-perms"`

	PathPerms os.FileMode
	FilePerms os.FileMode
}

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

func (cfg Config) GetEnv(env string) EnvDirSection {
	if conf, ok := cfg.Sections[env]; ok {
		conf.ParsePerms()
		return *conf
	}
	panic(fmt.Sprintf("envdir '%s' not found in config", env))
}

func ParsePerms(perms string, def os.FileMode) os.FileMode {
	if len(perms) > 0 {
		p, err := strconv.ParseUint(perms, 8, 32)
		if err != nil {
			panic(err)
		}
		return os.FileMode(p)
	}
	return def
}

func (e *EnvDirSection) ParsePerms() {
	e.PathPerms = ParsePerms(e.PathPermsString, DefaultPathPerms)
	e.FilePerms = ParsePerms(e.FilePermsString, DefaultFilePerms)
}
