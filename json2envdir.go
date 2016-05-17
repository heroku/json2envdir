package json2envdir

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/heroku/json2envdir/config"
	"github.com/satori/go.uuid"
)

var (
	Name string
	Env  map[string]interface{}
	JSON map[string]interface{}
)

// Funcs is a set of functions that can be used in
// the template value of an env var
type Funcs struct {
}

func (f Funcs) UUID() string {
	u := uuid.NewV4()
	return u.String()
}

func Parse(cfg config.Config, rawJSON string) error {
	json.Unmarshal([]byte(rawJSON), &JSON)

	Name = JSON["name"].(string)
	Env = JSON["env"].(map[string]interface{})

	envCfg := cfg.GetEnv(Name)

	os.MkdirAll(envCfg.Path, envCfg.PathPerms)

	funcs := Funcs{}
	for key := range Env {
		value := fmt.Sprintf("%v", Env[key])
		tmpl, err := template.New("").Parse(value)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(filepath.Join(envCfg.Path, key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, envCfg.FilePerms)
		if err != nil {
			return err
		}

		err = tmpl.Execute(f, funcs)
		if err != nil {
			return err
		}

		f.Close()
	}

	return nil
}
