#!/usr/bin/env bash
# Builds and runs the complete programs (code blocks with "package main") of the documentation and of the skill
# for AI agents against the code of the repository, so they don't describe the API that doesn't exist,
# and compares their outputs with the text blocks that follow them.
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
	# a program is a go block with "package main", the text block after it, if any, is its expected output
	n=$(awk -v dir="$work" -v n="$n" -v src="$file" '
		/^```go$/ { block = ""; inside = "go"; next }
		/^```text$/ { block = ""; inside = "text"; next }
		/^```$/ && inside {
			if (inside == "go" && block ~ /^package main/) {
				n++
				path = dir "/p" n
				system("mkdir -p " path)
				printf "%s", block > (path "/main.go")
				close(path "/main.go")
				print src > (path "/source")
				close(path "/source")
				expect = 1
			} else if (inside == "text" && expect) {
				printf "%s", block > (dir "/p" n "/expected")
				close(dir "/p" n "/expected")
				expect = 0
			} else if (inside == "go") {
				expect = 0
			}
			inside = ""
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

require (
	github.com/insei/valigo v0.1.0
	github.com/insei/valigo/huma v0.1.0
	github.com/insei/valigo/uuid v0.1.0
)

replace (
	github.com/insei/valigo => $root
	github.com/insei/valigo/huma => $root/huma
	github.com/insei/valigo/uuid => $root/uuid
)
MOD
cat "$root/go.sum" "$root/uuid/go.sum" "$root/huma/go.sum" | sort -u > "$work/go.sum"
cd "$work"
go mod tidy >/dev/null
go vet ./...
for dir in p*/; do
	if ! output=$(go run "./$dir"); then
		echo "the program from $(cat "$dir/source") failed" >&2
		exit 1
	fi
	if [ -f "$dir/expected" ] && [ "$output" != "$(cat "$dir/expected")" ]; then
		echo "the output of the program from $(cat "$dir/source") differs from the documented one:" >&2
		diff <(cat "$dir/expected") <(echo "$output") >&2 || true
		exit 1
	fi
done
echo "$n programs are built and run, their outputs match the documentation"
