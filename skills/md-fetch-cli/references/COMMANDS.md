# md-fetch command reference

## CLI fetch

```bash
md-fetch https://example.com
md-fetch --browser firefox https://example.com
md-fetch --browser curl https://example.com
```

## Save output

```bash
md-fetch --save https://example.com
md-fetch --save --filename output.md https://example.com
```

## Server mode

```bash
md-fetch serve
md-fetch serve --port 9090
```

## API call

```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://example.com"],"browser":"chrome"}'
```

## Tests

```bash
go test ./...
```
