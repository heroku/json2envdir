package json2envdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heroku/json2envdir/config"
	"github.com/satori/go.uuid"
)

func TestParse(t *testing.T) {
	dir, err := ioutil.TempDir("", "json2envdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	os.Setenv("FOO", "BAR")
	defer os.Setenv("FOO", "")

	cfg := config.Config{
		Sections: map[string]*config.EnvDirSection{
			"myapp": &config.EnvDirSection{
				Path: dir,
			},
		},
	}
	err = Parse(cfg, `{
		"name": "myapp",
		"env": {
			"MYVAR": "123",
			"MYVAR_TMPL": "{{.UUID}} 123",
			"MYVAR_ENV": "{{.Env \"FOO\"}} 123"
		}
	}`)
	if err != nil {
		t.Fatalf("Unexpected error in Parse: %s", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(dir, "MYVAR"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "123" {
		t.Fatalf("Value of $MYVAR should equal to 123, but it's %s", b)
	}

	b, err = ioutil.ReadFile(filepath.Join(dir, "MYVAR_TMPL"))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(string(b), " ")
	_, err = uuid.FromString(parts[0])
	if err != nil {
		t.Fatalf("Value of $MYVAR_TMPL should contain a UUID")
	}
	if parts[1] != "123" {
		t.Fatalf("Value of $MYVAR_TMPL should contain 123, but it's %s", b)
	}

	b, err = ioutil.ReadFile(filepath.Join(dir, "MYVAR_ENV"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "BAR 123" {
		t.Fatalf("Value of $MYVAR_ENV should equal to BAR 123, but it's %s", b)
	}
}
