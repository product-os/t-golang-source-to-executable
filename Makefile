IMAGENAME = golang-source-to-executable:latest

.PHONY: dockerize
dockerize:
	docker build -q -t ${IMAGENAME} .

test: OUTDIR=./test/out
test: dockerize
	@$(RM) -rf ./test/input.json
	@$(RM) -rf ${OUTDIR}
	@mkdir ${OUTDIR}
	@yq r -j ./test/artifact/balena.yml | jq '{input:{contract:.,artifactPath:"artifact"}}' >./test/input.json
	@docker run --rm -it \
		--env=INPUT=/input/input.json \
		--mount=type=bind,source=${PWD}/test,target=/input \
		--env=OUTPUT=/output/results.json \
		--mount=type=bind,source=${PWD}/test/out,target=/output \
		${IMAGENAME}
	@test -f ${OUTDIR}/results.json
	@jq '.' ${OUTDIR}/results.json
	@jq -r '.results[].artifactPath' ${OUTDIR}/results.json | xargs -I{} test -d ${OUTDIR}/{}
