export CGOENABLED=0
export GOOS=linux
export VERSION=`cat version`

build:
	@go build -o build/kube-localapi ./cmd/kube-localapi
	@go build -o build/kube-oauth2server ./cmd/kube-oauth2server
	@docker build -t coldog/kube-oauth2:$(VERSION) .

push:
	@docker push coldog/kube-oauth2:$(VERSION)

.PHONY: build push
