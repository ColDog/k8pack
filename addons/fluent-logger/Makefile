VERSION=`cat version`

build:
	@docker build -t coldog/fluent-logger:$(VERSION) .

push:
	@docker push coldog/fluent-logger:$(VERSION)

.PHONY: build push
