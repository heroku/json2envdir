package json2envdir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/heroku/json2envdir/config"
)

var (
	Name string
	Env  map[string]interface{}
	JSON map[string]interface{}
)

func Parse(cfg config.Config, rawJSON string) error {
	json.Unmarshal([]byte(rawJSON), &JSON)

	Name = JSON["name"].(string)
	Env = JSON["env"].(map[string]interface{})

	envCfg := cfg.GetEnv(Name)

	os.MkdirAll(envCfg.Path, envCfg.PathPerms)

	for key := range Env {
		value := fmt.Sprintf("%v", Env[key])
		ioutil.WriteFile(filepath.Join(envCfg.Path, key), []byte(value), envCfg.FilePerms)
	}

	return nil
}
