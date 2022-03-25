#!/bin/bash

set -e

DEBUG="${DEBUG:-1}"
SUITE=${SUITE:?missing env}
IMAGENAME=${IMAGENAME:?missing env}
OUTDIR=$(mktemp -d --tmpdir "tf-golang-${SUITE}"-XXXX)

jq -r '.' "./test/${SUITE}/input/input-contract.json"

artifactPath=$(jq -r '.results[0].artifactPath' "./test/${SUITE}/output/output-manifest.json")

# run the tf container on $SUITE
docker run --rm \
	--cap-add=CAP_SYS_ADMIN \
	--mount=type=bind,target=/input,source="${PWD}/test/${SUITE}/input",ro \
	--env=INPUT=/input/input-manifest.json \
	--mount=type=bind,target=/output,source="${OUTDIR}" \
	--env=OUTPUT=/output/output-manifest.json \
	--mount=type=tmpfs,target=/cache \
	--env=DEBUG="${DEBUG}" \
	"${IMAGENAME}"

# check an output manifest exists
test -f "${OUTDIR}"/output-manifest.json || {
	echo FAIL: missing output manifest
	exit 1
}

diff -q <(jq --sort-keys '.' ./test/"${SUITE}"/output/output-manifest.json) <(jq --sort-keys '.' "${OUTDIR}"/output-manifest.json) || {
	diff -u <(jq --sort-keys '.' ./test/"${SUITE}"/output/output-manifest.json) <(jq --sort-keys '.' "${OUTDIR}"/output-manifest.json) | bat -p
	echo FAIL: output manifest not matching
	exit 1
}

# print output manifest
jq '.' "${OUTDIR}"/output-manifest.json

# check the output manifest points to the correct artifact path
test "${artifactPath}" = "$(jq -r '.results[0].artifactPath' "${OUTDIR}"/output-manifest.json)" || {
	echo FAIL: artifactPath not matching
	exit 1
}

# check the artifact path exists
test -d "${OUTDIR}/${artifactPath}" || {
	echo FAIL: artifactPath missing
	exit 1
}

# test that all expected output files exist
# NOTE: this doesn't actually check they are identical
find "./test/${SUITE}/output/${artifactPath}/" -print0 |
	xargs -0 basename |
	find "${OUTDIR}/${artifactPath}/" -exec test -e {} ';'

echo "PASS: ${SUITE}"
