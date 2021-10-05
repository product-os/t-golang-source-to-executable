#!/bin/bash

set -e

INPUT_ARTIFACT_PATH="$(dirname "${INPUT}")/$(jq -r '.input.artifactPath' "${INPUT}")"
test -d "${INPUT_ARTIFACT_PATH}"

# TODO handle array?
INPUT_PLATFORM=$(jq -r '.input.contract.data.platforms[0]' "${INPUT}")
INPUT_TAGS=$(jq -r '.input.contract.data.tags' "${INPUT}")
INPUT_VERSION=$(jq -r '.input.contract.version' "${INPUT}")
INPUT_NAME=$(jq -r '.input.contract.name' "${INPUT}")

OUTPUT_TYPE=type-product-os-executable@1.0.1
OUTPUT_ARTIFACT_DIRNAME=artifacts
OUTPUT_VERSION_VARIABLE=version.Version
OUTPUT_FILENAME=$INPUT_NAME
OUTPUT_PLATFORM=$INPUT_PLATFORM
OUTPUT_VERSION=$INPUT_VERSION
OUTPUT_ARTIFACT_PATH="$(dirname "${OUTPUT}")/${OUTPUT_ARTIFACT_DIRNAME}"

RESULTS='[]'

function goBuild() {
	cd "${INPUT_ARTIFACT_PATH}" || exit 1

	go_build_args=

	# handle non-modules
	if ! [ -f go.mod ]; then
		mkdir -p "${GOPATH}/src/"
		ln -sfv "${INPUT_ARTIFACT_PATH}" "${GOPATH}/src/${INPUT_NAME}"
		cd "${GOPATH}/src/${INPUT_NAME}"
		export GO111MODULE=off
	fi

	if [ -n "${INPUT_TAGS}" ]; then
		go_build_args+="-tags ${INPUT_TAGS}"
	fi

	# TODO revision? build time stamp? package?
	(
		set -x
		go build -x \
			-ldflags "-X ${OUTPUT_VERSION_VARIABLE}=${INPUT_VERSION}" \
			-o "${OUTPUT_ARTIFACT_PATH}/${OUTPUT_FILENAME}" \
			${go_build_args} \
			./cmd/"${OUTPUT_FILENAME}"
	)

	RESULTS=$(json_append "${RESULTS}" "$(goBuildResult)")
}

function json_append() {
	jq -n -c --argjson total "${1}" --argjson new "${2}" '$total + [$new]'
}

function goBuildResult() {
	jq -n -c \
		--arg type "${OUTPUT_TYPE}" \
		--arg platform "${OUTPUT_PLATFORM}" \
		--arg filename "${OUTPUT_FILENAME}" \
		--arg version "${OUTPUT_VERSION}" \
		--arg artifactPath "${OUTPUT_ARTIFACT_DIRNAME}" \
		'{contract:{type:$type,data:{version:$version,platform:$platform,filename:$filename}},artifactPath:$artifactPath}'
}

function writeResultsJson() {
	jq -n -c --argjson results "${RESULTS}" '{results:$results}' >"${OUTPUT}"
}

mkdir -p "${OUTPUT_ARTIFACT_PATH}"

goBuild
writeResultsJson

# done
