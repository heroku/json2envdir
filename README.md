# JSON to envdir tool

In goes json, out goes envdir.

## JSON Format

```json
{
  "name: "myapp",
  "env": {
    "MYVAR": "123", 
    ...
  }
}
```

## Config File Format

**Default Config File**: `/etc/json2env.conf`

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
