BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-calendar: build-calendar
	$(BIN) -config ./configs/calendar-config.yml

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/calendar/Dockerfile .

run-calendar-img: build-calendar-img
	docker run $(DOCKER_IMG)

version: build-calendar
	$(BIN) version

test:
	go test -race ./internal/logger ./internal/storage/memory ./internal/server/http

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

down:
	docker-compose -f ./deployments/docker-compose.yml down

up: 
	docker-compose -f ./deployments/docker-compose.yml up -d --build

restart: down up

install-oapigen-deps:
	(which oapi-codegen > /dev/null) || go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

target: install-oapigen-deps
	go generate ./internal/server/http

build-scheduler:
	go build -v -o ./bin/scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler

run-scheduler: build-scheduler
	./bin/scheduler -config ./configs/scheduler-config.yml

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t scheduler:develop \
		-f build/scheduler/Dockerfile .

build-storer:
	go build -v -o ./bin/storer -ldflags "$(LDFLAGS)" ./cmd/storer

run-storer: build-storer
	./bin/storer -config ./configs/storer-config.yml

build-storer-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t storer:develop \
		-f build/storer/Dockerfile .

run-storer-img: build-storer-img
	docker run storer:develop

build: build-calendar build-scheduler build-storer

integration-tests:
	set -e ;\
	docker-compose -f ./integrationtests/docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./integrationtests/docker-compose.test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./integrationtests/docker-compose.test.yml down ;\
	exit $$test_status_code ;

integration-tests-cleanup:
	docker-compose -f ./integrationtests/docker-compose.test.yml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60;

prune:
	docker image prune -f
	docker volume prune -f

.PHONY: build run build-img run-img version test lint comp up down restart
