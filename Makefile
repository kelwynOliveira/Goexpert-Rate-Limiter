include .env
export

.PHONY: test
test: ## Run unit-tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html

.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down


## ---------- SCENARIOS
define scenario_ip
	echo "Running IP rate limit scenario..."; \
	for i in {1..4}; do \
		echo  "Request $$i: "; \
		curl -is -w "%{http_code} \n" -o /dev/null http://localhost:8080/api/v1/zipcode/03031040; \
	done; \
	echo "\nWait for block duration: $(BLOCK_DURATION)s"; \
	sleep $(BLOCK_DURATION); \
	echo "Request after block: "; \
	curl -is -w "%{http_code} \n" -o /dev/null http://localhost:8080/api/v1/zipcode/03031040 
endef

define scenario_token
	echo "Running token rate limit scenario..."; \
	for i in {1..6}; do \
		echo "Request $$i: "; \
		curl -is -w "%{http_code} \n" -o /dev/null -H "API_KEY: my-token" http://localhost:8080/api/v1/zipcode/03031040; \
	done; \
	echo "\nWait for block duration: $(BLOCK_DURATION)s"; \
	sleep $(BLOCK_DURATION); \
	echo "Request after block: "; \
	curl -is -w "%{http_code} \n" -o /dev/null -H "API_KEY: my-token" http://localhost:8080/api/v1/zipcode/03031040
endef

.PHONY: run
run: ## Run test scenarios
	@if [ "$(SCENARIO)" = "ip" ]; then \
		$(call scenario_ip); \
	elif [ "$(SCENARIO)" = "token" ]; then \
		$(call scenario_token); \
	elif [ "$(SCENARIO)" = "all" ]; then \
		$(call scenario_ip); \
		echo -e "\n----------------------------------------"; \
		$(call scenario_token); \
	else \
		echo "Please specify a valid SCENARIO: (ip, token, all)"; \
	fi
