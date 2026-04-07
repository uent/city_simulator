## 1. Update Go Version

- [x] 1.1 In `go.mod`, change `go 1.24.1` to `go 1.26.1`

## 2. Verify

- [x] 2.1 Run `go mod tidy` to refresh `go.sum` if needed
- [x] 2.2 Run `go build ./...` and confirm zero errors
- [x] 2.3 Run `go vet ./...` and confirm zero warnings
