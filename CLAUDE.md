# Go Project Guidelines

## Build Commands
- Initialize module (if not done): `go mod init learngo && go mod tidy`
- Build: `go build`
- Run: `go run .`
- Install: `go install`

## Test Commands
- Run all tests: `go test ./...`
- Run specific test: `go test -run TestName`
- Run tests with coverage: `go test -cover ./...`

## Lint Commands
- Format code: `go fmt ./...`
- Run linter: `golint ./...` (install with `go install golang.org/x/lint/golint@latest`)
- Vet code: `go vet ./...`

## Code Style Guidelines
- Use standard Go formatting (enforced by `go fmt`)
- Group imports: standard library first, then third-party packages
- Variable/function names: use camelCase for private, PascalCase for exported
- Comment all exported functions, variables, and types
- Error handling: check all errors, don't use panic/recover in normal operations
- Prefer early returns over nested if statements
- Use meaningful variable names that describe their purpose
- Organize code in related functions and packages