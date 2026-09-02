# Contributing

Thanks for helping with Pusula Serve.

## One product

The Go CLI (`cmd/pusula-serve`), the HTTP API (`internal/httpapi`), and the GitHub Pages console (`docs/`) are the same planner. Do not add a second calculator, a second GPU table, or a second set of product links.

Canonical links:

- Console: https://pusulainfra.github.io/pusula-serve/
- Source: https://github.com/PusulaInfra/pusula-serve
- X: https://x.com/pusulainfra

## Process

1. Fork the repository.
2. Create a branch (`git checkout -b feat/short-name`).
3. Keep engine numbers honest — no fake throughput, no invented cloud quotes.
4. Run `go test ./...` and `go build ./cmd/pusula-serve`.
5. Open a pull request against `main`.

## Code

- Go 1.22+.
- Console math in `docs/engine.js` should stay aligned with `internal/serve/analyze.go`.
- Do not commit secrets, GPU vendor credentials, or private model weights.
