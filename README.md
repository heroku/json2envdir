# JSON to envdir tool

In goes json, out goes envdir.

## JSON Format

```json
{
  "name": "myapp",
  "env": {
    "MYVAR": "123",
    "MYVAR_TMPL": "{{.UUID}} 123",
    ...
  }
}
```

Notice that you can have dynamic value for an env var. The format is in Go template. Currently only `{{.UUID}}` is supported.

## Config File Format

**Default Config File**: `/etc/json2envdir.conf`

```
[envdir "webapp"]
path = /www/webapp/env.d

[envdir "myapp"]
path = /etc/myapp/env.d
file-perms = 0640
path-perms = 0750
```

## Usage

```bash
json2envdir --file=app-env.json
```
