# aels

bcrypt based license server

kind of neat

## Examples

form element "license"

```
curl -v -d license=243261243130245343464765543338357169354434395a566f3976612e6a694d61426b4c4647364e4a46494c342f55792f624a67797545586b2e3047 localhost:8080
```

json payload key "license"

```
curl  -v -H "Content-Type:application/json" -d '{"license": "243261243130245343464765543338357169354434395a566f3976612e6a694d61426b4c4647364e4a46494c342f55792f624a67797545586b2e3047"}' localhost:8080
```
