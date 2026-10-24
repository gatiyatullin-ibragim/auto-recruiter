.PHONY: test test-race build frontend security

test:
	go test ./...

test-race:
	go test -race -cover ./...

build:
	go build ./...

frontend:
	cd web && npm run build

security:
	go vet ./...
	cd web && npm audit --audit-level=high
