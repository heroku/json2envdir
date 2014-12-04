package json2envdir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var (
	FilePerms os.FileMode = 0640
	PathPerms os.FileMode = 0750
	Path      string
	Env       map[string]interface{}
	JSON      map[string]interface{}
)

func Parse(rawJSON string) error {
	json.Unmarshal([]byte(rawJSON), &JSON)

	Path = JSON["path"].(string)
	Env = JSON["env"].(map[string]interface{})
	if str, ok := JSON["path-perms"]; ok {
		perm, err := strconv.ParseUint(str.(string), 8, 32)
		if err != nil {
			return err
		}
		PathPerms = os.FileMode(perm)
	}
	if str, ok := JSON["file-perms"]; ok {
		perm, err := strconv.ParseUint(str.(string), 8, 32)
		if err != nil {
			return err
		}
		FilePerms = os.FileMode(perm)
	}

	os.MkdirAll(Path, PathPerms)

	for key := range Env {
		value := fmt.Sprintf("%v", Env[key])
		ioutil.WriteFile(filepath.Join(Path, key), []byte(value), FilePerms)
	}

	return nil
}
