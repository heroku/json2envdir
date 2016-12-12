package json2envdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/json2envdir/config"
	"github.com/satori/go.uuid"
)

func TestParseNoConfig(t *testing.T) {
	testDir, err := ioutil.TempDir("", "json2envdirTest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(testDir) }()

	cfg := config.Config{
		Sections: map[string]*config.EnvDirSection{
			"cheetah": &config.EnvDirSection{
				Path: []string{testDir},
			},
		},
	}

	err = Process(cfg, `{
		"name": "not-cheetah",
		"env": {
			"VAR": "{{.UUID}}"
		}
	}`)

	if err != nil {
		t.Fatalf("Unexpected failure: %q", err)
	}
}

func TestParse(t *testing.T) {
	dir1, err := ioutil.TempDir("", "json2envdir1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir("", "json2envdir2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir2)

	cfg := config.Config{
		Sections: map[string]*config.EnvDirSection{
			"myapp": &config.EnvDirSection{
				Path: []string{dir1, dir2},
			},
		},
	}
	err = Process(cfg, `{
		"name": "myapp",
		"env": {
			"MYVAR": "123",
			"MYVAR_TMPL": "{{.UUID}}"
		}
	}`)
	if err != nil {
		t.Fatalf("Unexpected error in Parse: %s", err)
	}

	b, err := ioutil.ReadFile(filepath.Join(dir1, "MYVAR"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "123" {
		t.Fatalf("Value of $MYVAR in dir1 should equal to 123, but it's %s", b)
	}

	b, err = ioutil.ReadFile(filepath.Join(dir1, "MYVAR_TMPL"))
	if err != nil {
		t.Fatal(err)
	}
	uuid1, err := uuid.FromString(string(b))
	if err != nil {
		t.Fatalf("Value of $MYVAR_TMPL in dir1 should contain a UUID")
	}

	b, err = ioutil.ReadFile(filepath.Join(dir2, "MYVAR"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "123" {
		t.Fatalf("Value of $MYVAR in dir2 should equal to 123, but it's %s", b)
	}

	b, err = ioutil.ReadFile(filepath.Join(dir2, "MYVAR_TMPL"))
	if err != nil {
		t.Fatal(err)
	}
	uuid2, err := uuid.FromString(string(b))
	if err != nil {
		t.Fatalf("Value of $MYVAR_TMPL in dir2 should contain a UUID")
	}

	if uuid1 != uuid2 {
		t.Fatalf("uuid1 should equal to uuid2: uuid1=%s uuid2=%s", uuid1, uuid2)
	}
}
