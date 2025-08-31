# Continuous Integration and Coverage

This guide shows how to run Ticketr’s CI and report coverage with Codecov.

## GitHub Actions

Minimal workflow that vets, lints, tests with coverage, and uploads to Codecov:

```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
          args: --timeout=5m

      - name: Test
        run: go test -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt ./...

      - name: Upload coverage artifact
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: coverage
          path: coverage.txt

      - name: Upload coverage reports to Codecov
        if: success() || failure()
        uses: codecov/codecov-action@v5
        with:
          token: ${{ secrets.CODECOV_TOKEN }}
          slug: karolswdev/ticketr
```

Notes:
- `-covermode=atomic` pairs well with `-race`.
- `-coverpkg=./...` ensures a single `coverage.txt` at the repo root.
- Codecov v5 autodiscovers `coverage.txt`; explicit `files:` is optional.

## Local Checks

Run the same steps locally:

```bash
go vet ./...
golangci-lint run
go test -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt ./...
```

Open an HTML coverage report:

```bash
go tool cover -html=coverage.txt -o coverage.html
```

