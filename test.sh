#!/bin/sh

set -e

echo "Running go vet..."
go vet ./src/

echo ""
echo "Checking for REMOVE ME..."

grep -RIn "REMOVE ME" src/ || true

echo ""
echo "Checking for TODOs..."

# matches:
# TODO: something
# todo: something
# todo something
grep -RInE '([Tt]ODO:|[Tt]odo:|[Tt]odo ).*' src/ || true
