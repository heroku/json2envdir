package json2envdir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
	funcs := Funcs{}
	for key := range j.Env {
		value := fmt.Sprintf("%v", j.Env[key])

		var result bytes.Buffer
		tmpl, err := template.New("").Parse(value)
		if err != nil {
			return err
		}
		err = tmpl.Execute(&result, funcs)
		if err != nil {
			return err
		}

		for _, path := range envCfg.Path {
			err = os.MkdirAll(path, envCfg.PathPerms)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(filepath.Join(path, key), result.Bytes(), envCfg.PathPerms)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
