package json2envdir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	FilePerms os.FileMode = 0644
	PathPerms os.FileMode = 0755
	Path      string
	Env       map[string]interface{}
	JSON      map[string]interface{}
)

func Parse(rawJSON string) error {
	json.Unmarshal([]byte(rawJSON), &JSON)

	Path = JSON["path"].(string)
	Env = JSON["env"].(map[string]interface{})
	if perm, ok := JSON["path-perms"]; ok {
		PathPerms = perm.(os.FileMode)
	}
	if perm, ok := JSON["file-perms"]; ok {
		FilePerms = perm.(os.FileMode)
	}

	os.MkdirAll(Path, PathPerms)

	for key := range Env {
		value := fmt.Sprintf("%v", Env[key])
		ioutil.WriteFile(filepath.Join(Path, key), []byte(value), FilePerms)
	}

	return nil
}
