MODE ?= build
IMAGENAME = golang-source-to-executable:latest

.PHONY: dockerize
dockerize:
	docker build -t ${IMAGENAME} .

test: test-module test-legacy

test-module: dockerize
	@env SUITE=$@ IMAGENAME=${IMAGENAME} ./test/run.sh

test-legacy: dockerize
	@env SUITE=$@ IMAGENAME=${IMAGENAME} ./test/run.sh
