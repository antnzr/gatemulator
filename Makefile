# Definitions
ROOT                    := $(PWD)
GOLANG_DOCKER_IMAGE     := golang:1.22-alpine
GOLANG_DOCKER_CONTAINER := gatemulator

dev:
	@CompileDaemon -exclude-dir=".git,migrations" \
		-command="./bin/app" \
		-build="go build -o ./bin/app cmd/gatemulator/main.go" \
		-color -log-prefix=false

build:
	@go build -o ./bin/app cmd/gatemulator/main.go

run:
	@./bin/app

start: build run

build_docker:
	@docker build -t gatemulator .

run_docker:
	@docker run -it --rm -v ./gatemulator.db:/app/gatemulator.db \
		-p 34000:34000 \
		--name gatemulator gatemulator

load_test:
	@docker run --rm -it -v ${PWD}/test/load/:/scripts \
  	--name artillery artilleryio/artillery:latest \
  	run /scripts/artillery.yml
