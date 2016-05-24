package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfgStr := `[envdir "myapp"]
path = path1
path = path2
file-perms = 0640
path-perms = 0750
`

	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "example.conf")
	err = ioutil.WriteFile(path, []byte(cfgStr), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	cfg := LoadConfig(path)
	section := cfg.Sections["myapp"]
	if section == nil {
		t.Fatalf("no myapp section is found")
	}

	if !reflect.DeepEqual(section.Path, []string{"path1", "path2"}) {
		t.Fatalf("section should contain path1 and path2: %v", section.Path)
	}

	if section.PathPermsString != "0750" {
		t.Fatalf("path perms should equal to 0750", section.PathPermsString)
	}

	if section.FilePermsString != "0640" {
		t.Fatalf("file perms should equal to 0640", section.FilePermsString)
	}
}
