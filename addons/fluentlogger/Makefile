version = `cat version`

build-agent:
	@docker build -t coldog/fluentlogger-agent:$(version) ./agent
	@docker tag coldog/fluentlogger-agent:$(version) coldog/fluentlogger-agent:latest
.PHONY: build-agent

build-collector:
	@docker build -t coldog/fluentlogger-collector:$(version) ./collector
	@docker tag coldog/fluentlogger-collector:$(version) coldog/fluentlogger-collector:latest
.PHONY: build-collector

push:
	@docker push coldog/fluentlogger-collector:$(version)
	@docker push coldog/fluentlogger-agent:$(version)
.PHONY: push

build: build-collector build-agent
.PHONY: build
