# JSON to envdir tool

In goes json, out goes envdir.

## JSON Format
testing

```json
{
  "name": "myapp",
  "env": {
    "MYVAR": "123",
    "MYVAR_TMPL": "{{.UUID}} 123",
    "MYVAR_HEX": "{{.HEX 32}} 123",
    ...
  }
}
```

Notice that you can have dynamic value for an env var. The format is in Go template. Currently `{{.UUID}}` and `{{.HEX N}}` are supported. HEX accepts an integer for length - the resulting string will be double this length.

## Config File Format

**Default Config File**: `/etc/json2envdir.conf`

```ini
[envdir "webapp"]
path = /www/webapp/env.d

[envdir "myapp"]
path = /etc/myapp1/env.d
path = /etc/myapp2/env.d
file-perms = 0640
path-perms = 0750
```

## Usage

```bash
json2envdir --file=app-env.json
```
