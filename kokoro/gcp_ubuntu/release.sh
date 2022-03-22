#!/bin/bash

# Fail on any error.
set -e

# Display commands being run.
# WARNING: please only enable 'set -x' if necessary for debugging, and be very
#  careful if you handle credentials (e.g. from Keystore) with 'set -x':
#  statements like "export VAR=$(cat /tmp/keystore/credentials)" will result in
#  the credentials being printed in build logs.
#  Additionally, recursive invocation with credentials as command-line
#  parameters, will print the full command, with credentials, in the build logs.
# set -x

# Code under repo is checked out to ${KOKORO_ARTIFACTS_DIR}/github.
# The final directory name in this path is determined by the scm name specified
# in the job configuration.
cd "${KOKORO_ARTIFACTS_DIR}/github/kpt/porch"

# Run second script in docker image containing dependencies
docker run --rm -i \
  --volume "${KOKORO_ARTIFACTS_DIR}:${KOKORO_ARTIFACTS_DIR}" \
  --workdir "${KOKORO_ARTIFACTS_DIR}/github/kpt/porch" \
  --entrypoint "${KOKORO_ARTIFACTS_DIR}/github/kpt/porch/build.sh" \
  "golang:1.17.8"


IMAGE_REPO=gcr.io/sdowell-catalyst-dev IMAGE_TAG="${KOKORO_GIT_COMMIT}" make build-images


TAG_NAME="test-foo"
REF="HEAD"

git tag "${TAG_NAME}" "${REF}"

git push origin "${TAG_NAME}"
