# Project Completion Checklist

This checklist tracks implemented requirements according to the final project specification.

---

# Part 1 — MVP (25 points)

| Done | Requirement |
| :--: | ----------- |
| [x] | `POST /api/shorten` returns **201** with `code` and `short_url` |
| [x] | **Idempotency:** second `POST` with the same URL returns the same code |
| [x] | `GET /{code}` returns **302** with correct `Location` |
| [x] | Unknown code → **404**; bad/missing URL → **400** |
| [x] | URL validation and no server-side fetch of long URL |
| [x] | Codes are 6–8 characters; collision handling implemented |
| [x] | `-base` flag used for `short_url` |
| [x] | `httptest`: shorten + redirect; table tests for invalid URL and unknown code |
| [x] | Concurrent duplicate shorten test; `go test -race ./...` passes |

---

# Part 2 — API & Errors (25 points)

| Done | Requirement |
| :--: | ----------- |
| [x] | Metadata route **200 / 404** with correct JSON |
| [x] | `ErrNotFound`, `ErrInvalidURL` implemented |
| [x] | `%w` wrapping + `errors.Is` used in HTTP mapping |
| [x] | `Store` interface + fake/mock used in tests |
| [x] | Tests for metadata route and error mapping |
| [x] | Idempotency still works through Store layer |

---

# Part 3 — Performance & Measurement (25 points)

| Done | Requirement |
| :--: | ----------- |
| [x] | Server timeouts configured |
| [x] | Mutex choice documented in `DECISIONS.md` |
| [x] | Benchmarks for shorten and redirect paths |
| [x] | README includes benchmark line and profiling insight |
| [x] | Tests and `-race` remain green |

---

# Part 4 — Persistence (25 points)

| Done | Requirement |
| :--: | ----------- |
| [x] | Persistent Store implemented with PostgreSQL + GORM |
| [x] | Startup loads existing links |
| [x] | Data persisted before returning success response |
| [x] | Restart simulation test implemented |
| [x] | Configuration supports memory vs persistent storage |
| [x] | `go test -race ./...` passes with persistence |

---

# Part 5 — Millions of Requests (Optional Bonus)

| Done | Requirement |
| :--: | ----------- |
| [x] | `DECISIONS.md`: Load Balancer → multiple Go apps → shared store |
| [ ] | CDN / edge caching for redirects |
| [x] | Write-path scaling using rate limiting |
| [ ] | Sharding / partitioning strategy |
| [ ] | Bonus code: cache layer, load-test script, or additional measurements |

---

# Part 6 — Production Habits (Optional Bonus)

| Done | Requirement |
| :--: | ----------- |
| [x] | Graceful shutdown working and documented |
| [x] | Rate limit on create endpoint |
| [ ] | Domain policy in `DECISIONS.md` |
| [x] | Logging policy documented (what is logged and what is avoided) |
| [ ] | Bonus: metrics, pprof flag, or structured logging |

---

# Verification

Commands verified:

```bash
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go vet ./...
```


Current statement coverage:

80.4%