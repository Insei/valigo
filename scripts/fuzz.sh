#!/usr/bin/env bash
# Runs every fuzz target of the core for FUZZTIME (30s by default). A failing input is saved in testdata/fuzz
# of the package, commit it with the fix: go test runs it as a regression test.
set -euo pipefail

cd "$(dirname "$0")/.."
time=${FUZZTIME:-30s}
for pkg in ./rules/str ./shared ./internal/tags; do
	for target in $(go test -list '^Fuzz' "$pkg" | grep '^Fuzz'); do
		echo "$pkg $target"
		go test -run '^$' -fuzz "^$target\$" -fuzztime "$time" "$pkg"
	done
done
