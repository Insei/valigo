#!/usr/bin/env bash
# Checks that the modules of modules.sh require a released version of the core: the replace of their go.mod
# is ignored by the users, a pseudo-version of the core, v0.0.0-..., can't be downloaded.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
. "$root/scripts/modules.sh"
failed=0
for module in $MODULES; do
	version=$(cd "$root/$module" && go list -m -f '{{.Version}}' github.com/insei/valigo)
	if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
		echo "$module requires github.com/insei/valigo $version, a released version expected" >&2
		failed=1
	fi
done
exit $failed
