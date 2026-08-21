# Store Inspection Record Ledger

The ledger keeps store inspection material in a Bolt-backed local database. It supports deterministic CSV import and validation, reviewer decisions, linked record ordering, search and pagination, collaboration notes, reports, and archival.

## Run

```bash
go run ./cmd/store-ledger health
go run ./cmd/store-ledger import
go run ./cmd/store-api
```

The HTTP service exposes `/health`, `/batches/import`, `/records`, `/records/review`, `/records/note`, and `/records/publish`.

## Test

```bash
go build ./...
go test -count=1 ./...
```

The standard `go run` smoke commands are:

```bash
go run ./cmd/store-ledger health
go run ./cmd/store-api
```

The regression test `TestBusiness07Regression` intentionally reproduces the injected equal-score ordering defect for batch `RB2177-07`.
See `BUG_REPRO.md` for the captured regression output and architecture reproduction evidence.
