#!/bin/bash
set -euo pipefail

TOPLEVEL=$(git rev-parse --show-toplevel)
export GOBIN="${TOPLEVEL}/bin"

GITCOMMIT="$(git rev-parse HEAD)"
TAG="$(git describe --tags --match 'v*' --abbrev=0)"
DESCRIBE="$(git describe --tags --match 'v*' --abbrev=40 --dirty --broken)"
COMMITS_SINCE_TAG="$(git log --format=%H ${TAG}..HEAD | wc -l | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"


LDFLAGS=""
LDFLAGS+="-X=main.GitCommit=${GITCOMMIT}"
LDFLAGS+=" -X \"main.GitDescribe=${DESCRIBE}\""
LDFLAGS+=" -X \"main.GitTag=${TAG}\""
LDFLAGS+=" -X \"main.CommitsSinceTag=${COMMITS_SINCE_TAG}\""


go fmt ./...
go test ./...
set -x
go install -ldflags "${LDFLAGS}" ./cmd/...
