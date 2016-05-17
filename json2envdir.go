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

type JSON struct {
	Name string                 `json:"name"`
	Env  map[string]interface{} `json:"env"`
}

// Funcs is a set of functions that can be used in
// the template value of an env var
type Funcs struct {
}

func (f Funcs) UUID() string {
	u := uuid.NewV4()
	return u.String()
}

func Parse(cfg config.Config, rawJSON string) error {
	var j JSON
	err := json.Unmarshal([]byte(rawJSON), &j)
	if err != nil {
		return err
	}

	envCfg := cfg.GetEnv(j.Name)
	err = os.MkdirAll(envCfg.Path, envCfg.PathPerms)
	if err != nil {
		return err
	}

	funcs := Funcs{}
	for key := range j.Env {
		value := fmt.Sprintf("%v", j.Env[key])
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
