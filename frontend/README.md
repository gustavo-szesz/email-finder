# Frontend

This is the Angular UI for the Go email-finder API.

## Run it

Start the Go API from the repository root:

```bash
go run ./cmd/api/main.go
```

Start the Angular app from this folder:

```bash
npm start
```

The Angular dev server proxies `/analyze` to `http://localhost:8080`, so the UI can talk to the Go backend without extra CORS setup during local development.

## Build

```bash
npm run build
```
