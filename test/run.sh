#!/bin/sh

set -e

TEST=${TEST:?missing env}
IMAGENAME=${IMAGENAME:?missing env}
OUTDIR=$(mktemp -d --tmpdir "${TEST}"-XXXX)

yq r -j "./test/${TEST}/balena.yml" | jq --arg artifactPath "${TEST}" '{input:{contract:.,artifactPath:$artifactPath}}' >./test/input.json
jq '.' ./test/input.json

docker run --rm -it \
	--env=INPUT=/input/input.json \
	--mount=type=bind,source="${PWD}"/test,target=/input \
	--env=OUTPUT=/output/results.json \
	--mount=type=bind,source="${OUTDIR}",target=/output \
	"${IMAGENAME}"

test -f "${OUTDIR}"/results.json
jq '.' "${OUTDIR}"/results.json
jq '.results[].artifactPath' "${OUTDIR}"/results.json | xargs -I{} test -d "${OUTDIR}/{}"

echo "PASS: ${TEST}"
