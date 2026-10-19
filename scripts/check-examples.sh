#!/usr/bin/env bash
# Builds and runs the complete programs (code blocks with "package main") of the documentation and of the skill
# for AI agents against the code of the repository, so they don't describe the API that doesn't exist.
# Usage: scripts/check-examples.sh [file.md...], all docs and the skill by default.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
if [ $# -eq 0 ]; then
	set -- "$root"/docs/*.md "$root"/docs/*/*.md "$root"/skills/valigo/SKILL.md
fi
# EXAMPLES_WORKDIR keeps the built programs to run them by hand
if [ -n "${EXAMPLES_WORKDIR:-}" ]; then
	work=$EXAMPLES_WORKDIR
	mkdir -p "$work"
else
	work=$(mktemp -d)
	trap 'rm -rf "$work"' EXIT
fi

n=0
for file in "$@"; do
	n=$(awk -v dir="$work" -v n="$n" -v src="$file" '
		/^```go$/ { block = ""; inside = 1; next }
		/^```$/ && inside {
			inside = 0
			if (block ~ /^package main/) {
				n++
				path = dir "/p" n
				system("mkdir -p " path)
				printf "%s", block > (path "/main.go")
				close(path "/main.go")
				print src > (path "/source")
				close(path "/source")
			}
			next
		}
		inside { block = block $0 "\n" }
		END { print n }
	' "$file")
done

if [ "$n" -eq 0 ]; then
	echo "no programs in $*" >&2
	exit 1
fi
cat > "$work/go.mod" <<MOD
module examplescheck

go 1.27

require github.com/insei/valigo v0.0.0

replace github.com/insei/valigo => $root
MOD
cp "$root/go.sum" "$work/go.sum"
cd "$work"
go mod tidy >/dev/null
go vet ./...
for dir in p*/; do
	if ! go run "./$dir" >/dev/null; then
		echo "the program from $(cat "$dir/source") failed" >&2
		exit 1
	fi
done
echo "$n programs are built and run"
