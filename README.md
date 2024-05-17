# Server Side Events

## Usage

### 1. Start SSE server

```sh
go run .
```

### 2. Open static web page with SSE listener

```sh
firefox static/index.html
# or
firefox http://localhost:3000/static
```

### 3. Trigger SSE event

```sh
curl http://localhost:3000/ping
```

The html page should have been updated.
