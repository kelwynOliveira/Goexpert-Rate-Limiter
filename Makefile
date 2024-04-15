build:
	@go build -ldflags="-w -s" -o ./bin/api ./cmd/api/main.go

run:
	@air -c .air.toml

test:
	@./scripts/test.sh

test_k6_smoke:
	@echo "==> Running smoke test with k6 via Docker"
	@docker compose run k6 run /scripts/smoke/smoke.test.js

test_k6_stress:
	@echo "==> Running stress test with k6 via Docker (will take a while to complete)"
	@docker compose run k6 run /scripts/stress/stress.test.js

install-deps:
	go install github.com/cosmtrek/air@latest
	go mod tidy

setup: install-deps

clean:
	rm -rf ./bin ./tmp ./coverage.txt
