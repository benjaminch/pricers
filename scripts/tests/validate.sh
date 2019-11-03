#!/bin/bash

set -e

AUTOFMT=true
COVERAGE=true
REPORTCARD=false # FIXME: reactivate once support for Go 1.12
VET=true

# Build a list of all the top-level directories in the project.
for DIRECTORY in */ ; do
  GOGLOB="$GOGLOB ${DIRECTORY%/}"
done
GOGLOB="${GOGLOB/ bin/}"
GOGLOB="${GOGLOB/ docs/}"
GOGLOB="${GOGLOB/ vendor/}"

# Check that there are no formatting issues
if $AUTOFMT; then
  COMMAND="go fmt ./..."
  echo "Running: $COMMAND"
  `$COMMAND`
fi

if $VET; then
  # Fix for the go 1.10 vet bug (https://github.com/w0rp/ale/issues/1358)
  COMMAND="go vet ./..."
  echo "Running: $COMMAND"
  `$COMMAND`
fi

# Tests
echo "Running: Tests"
go test -race ./... -cover -covermode=atomic -coverprofile=coverage.out

# Report card generates a report on the quality of an open source go project.
if $REPORTCARD; then
  go get github.com/gojp/goreportcard
  echo "Running: GoReportCard"
  goreportcard-cli -v
fi