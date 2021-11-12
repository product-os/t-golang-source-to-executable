#!/bin/sh

set -e

SUITE=${SUITE:?missing env}
IMAGENAME=${IMAGENAME:?missing env}
OUTDIR=$(mktemp -d --tmpdir "tf-golang-${SUITE}"-XXXX)

jq -r '.' "./test/${SUITE}/input/input-contract.json"

artifactPath=$(jq -r '.input.artifactPath' "./test/${SUITE}/input/input-contract.json")

# run the tf container on $SUITE
docker run --rm -it \
	--mount=type=bind,source="${PWD}/test/${SUITE}/input",target=/input,ro \
	--env=INPUT=/input/input-contract.json \
	--mount=type=bind,source="${OUTDIR}",target=/output \
	--env=OUTPUT=/output/output-manifest.json \
	--mount=type=tmpfs,target=/cache \
	--env=DEBUG=1 \
	"${IMAGENAME}"

# check an output manifest exists
test -f "${OUTDIR}"/output-manifest.json || {
	echo missing output manifest
	exit 1
}

# print output manifest
jq '.' "${OUTDIR}"/output-manifest.json

# check the output manifest points to the correct artifact path
test "${artifactPath}" = "$(jq -r '.results[0].artifactPath' "${OUTDIR}"/output-manifest.json)" || {
	echo artifactPath not matching
	exit 1
}

# check the artifact path exists
test -d "${OUTDIR}/${artifactPath}" || {
	echo artifactPath missing
	exit 1
}

# test that all expected output files exist
# NOTE: this doesn't actually check they are identical
find "./test/${SUITE}/output/${artifactPath}/" -type f -exec basename {} ';' |
	find "${OUTDIR}/${artifactPath}/" -type f -exec test -f {} ';'

echo "PASS: ${SUITE}"
