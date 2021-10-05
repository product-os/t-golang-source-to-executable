IMAGENAME = golang-source-to-executable:latest

.PHONY: dockerize
dockerize:
	docker build -q -t ${IMAGENAME} .

test: test-module test-legacy

test-module: dockerize
	@env TEST=$@ IMAGENAME=${IMAGENAME} ./test/run.sh

test-legacy: dockerize
	@env TEST=$@ IMAGENAME=${IMAGENAME} ./test/run.sh
