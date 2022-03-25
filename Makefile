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

test-mvp: dockerize
	@env SUITE=$@ IMAGENAME=${IMAGENAME} ./test/run.sh

tester:
	docker build -t t-tester ../t-tester
	docker run --rm -it --privileged --name=tester \
		--mount=type=bind,src=${PWD},target=/input \
		--mount=type=tmpfs,target=/output \
		--mount=type=volume,src=dindcache,target=/home/rootless/.local/share/docker \
		t-tester
