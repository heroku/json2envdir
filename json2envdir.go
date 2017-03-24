package json2envdir

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/heroku/json2envdir/config"
	"github.com/satori/go.uuid"
)

type payload struct {
	Name string                 `json:"name"`
	Env  map[string]interface{} `json:"env"`
}

// Funcs is a set of functions that can be used in
// the template value of an env var
type Funcs struct {
}

// UUID is a function exposed to the templates for generating a UUIDv4
func (f Funcs) UUID() string {
	u := uuid.NewV4()
	return u.String()
}

// Hex allows templates to generate a hex of a specific length
func (f Funcs) HEX(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}

// Process parses and populates the envdirs specified in the config if matching
// data is found in the json payload.
func Process(cfg config.Config, rawJSON string) error {
	var p payload
	err := json.Unmarshal([]byte(rawJSON), &p)
	if err != nil {
		return err
	}

	envCfg, found := cfg.GetEnv(p.Name)
	if !found {
		fmt.Printf("Skipping unconfigured entry: %s", p.Name)
		return nil
	}

	funcs := Funcs{}
	for key := range p.Env {
		value := fmt.Sprintf("%v", p.Env[key])

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
