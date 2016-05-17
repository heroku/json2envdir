package json2envdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/json2envdir/config"
)

func TestParse(t *testing.T) {
	dir, err := ioutil.TempDir("", "json2envdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cfg := config.Config{
		Sections: map[string]*config.EnvDirSection{
			"myapp": &config.EnvDirSection{
				Path: dir,
			},
		},
	}
	err = Parse(cfg, `{ "name": "myapp", "env": { "MYVAR": "123" } }`)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(filepath.Join(dir, "MYVAR"))
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "123" {
		t.Fatalf("Value of $MYVAR should equal to 123, but it's %s", b)
	}
}
