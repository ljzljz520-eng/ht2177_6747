# BUG_REPRO

The injected defect is in the `RB2177-07` import-and-validate workflow. Equal-score records are ordered by identifier instead of preserving their import sequence. The same faulty order is returned on the immediate repeated import.

## Reproduction

### Go test

```bash
go test ./internal/service -run TestBusiness07Regression
go test -count=1 ./internal/service -run TestBusiness07Regression
```

- Command: `go test -count=1 ./...`
- Exit status: `1`

Observed failure:

```text
same-score records must retain import sequence, got 1
```

The failure is expected for the bugbase project and is retained as the regression signal.

## Architecture reproduction

### linux/arm64

- Go build: exit `0`
- Go test: exit `1` (the intentional `TestBusiness07Regression` failure above)
- Go run smoke ./cmd/store-ledger: exit `0`
- Go run smoke ./cmd/store-api: exit `0`

### linux/amd64

- Go build: exit `0`
- Go test: exit `1` (the intentional `TestBusiness07Regression` failure above)
- Go run smoke ./cmd/store-ledger: exit `0`
- Go run smoke ./cmd/store-api: exit `0`
