default:
	make tidy vet build
.PHONY: default

build: tidy
	go build -o _build/server ./cmd/server
.PHONY: build

tidy:
	go mod tidy
.PHONY: tidy

vet:
	go vet ./...
.PHONY: vet

run:
	go run ./cmd/server
.PHONY: run

# Docker
build-docker:
	docker compose build
.PHONY: build-docker

up:
	docker compose --profile mon-go up -d
.PHONY: up

down:
	docker compose --profile mon-go down
.PHONY: down

up-mongo:
	docker compose --profile mongo up -d
.PHONY: up-mongo

down-mongo:
	docker compose --profile mongo down
.PHONY: down-mongo

# Build images and run containers (build + up)
run-docker: build-docker up
.PHONY: run-docker

c:
	rm -rf _build
.PHONY: c
