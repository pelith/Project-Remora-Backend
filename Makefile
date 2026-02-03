.PHONY: build docker release lint gci-format test coverage gen-migration-sql sqlc abigen
NAME=remora

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o "${NAME}_$(app)" ./cmd/$(app)

app?=api
tag?=latest

docker:
	docker build \
	--build-arg APP=$(app) \
	-t $(NAME):$(tag) . 

TAG := v$(shell date -u '+%Y.%m.%d.%H.%M.%S')

release: 
	git tag $(TAG)
	git push origin $(TAG)

# Tools

lint:
	@golangci-lint run ./... -c ./.golangci.yml

gci-format:
	@gci write --skip-generated -s standard -s default -s "Prefix(remora)" ./

test:
	@go test ./... -race  

coverage:
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out

# SQL

DATETIME=$(shell date -u '+%Y%m%d%H%M%S')

gen-migration-sql:
	@( \
	printf "Enter file name: "; read -r FILE_NAME; \
	touch database/migrations/$(DATETIME)_$$FILE_NAME.up.sql; \
	touch database/migrations/$(DATETIME)_$$FILE_NAME.down.sql; \
	)

gen-seed-sql:
	@( \
	printf "Enter file name: "; read -r FILE_NAME; \
	printf "Enter env: "; read -r ENV; \
	mkdir -p database/seeds/$$ENV; \
	touch database/seeds/$$ENV/$(DATETIME)_$$FILE_NAME.up.sql; \
	touch database/seeds/$$ENV/$(DATETIME)_$$FILE_NAME.down.sql; \
	)

sqlc:
	sqlc generate -f ./database/sqlc.yml

sqlc-lint:
	sqlc vet -f ./database/sqlc.yml

# Contracts

abigen:
	@abigen --abi contracts/stateview/StateView.json --pkg contracts --type StateView --out internal/liquidity/repository/contracts/stateview.go
