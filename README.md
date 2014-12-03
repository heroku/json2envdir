# JSON to envdir tool

In goes json, out goes envdir.

## Format

```json
{
  "path": "/etc/myapp/env.d",
  "env": {
    "MYVAR": "123", 
    ...
  },
  "path-perms": 0755  // optional
  "file-perms": 0644  // optional
}
```
