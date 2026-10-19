#!/usr/bin/env bash
# Checks that the unit tests of every package cover 100% of its statements, in the root, uuid and huma modules.
# Every package is measured by its own tests only, so a package can't rely on the tests of another one.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
profile=$(mktemp)
trap 'rm -f "$profile"' EXIT

failed=0
for module in "$root" "$root/uuid" "$root/huma"; do
	cd "$module"
	go test -count=1 -coverprofile="$profile" ./... >/dev/null
	# the blocks with zero count, a block is printed as file:line.col,line.col
	uncovered=$(awk 'NR > 1 && $NF == 0 { print $1 }' "$profile")
	total=$(go tool cover -func="$profile" | awk '/^total:/ { print $NF }')
	echo "$(basename "$module"): $total"
	if [ -n "$uncovered" ]; then
		echo "uncovered blocks:" >&2
		echo "$uncovered" >&2
		failed=1
	fi
done
exit $failed
