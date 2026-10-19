#!/usr/bin/env bash
#
# Runs the benchmarks on the head of the pull request and on the target branch, compares them by benchstat
# and posts the summary to the pull request: one comment, updated after every push.
#
# It is used by .github/workflows/bench-diff.yml on pull_request events and expects the variables of GitHub Actions:
#   GITHUB_BASE_REF, GITHUB_REPOSITORY, GITHUB_EVENT_PATH, GITHUB_API_URL, GITHUB_TOKEN (pull-requests: write).
# Without GITHUB_TOKEN the report is only in the logs and in the summary of the job.

set -euo pipefail

COUNT="${BENCH_COUNT:-10}"
BENCH_PATTERN='^Benchmark'
BENCH_PKGS=(.)
MARKER='<!-- valigo-bench-diff -->'

: "${GITHUB_BASE_REF:?must run from pull_request event}"
: "${GITHUB_REPOSITORY:?}"
: "${GITHUB_EVENT_PATH:?}"

PR_NUMBER="$(jq -r '.pull_request.number // empty' "$GITHUB_EVENT_PATH")"
API_URL="${GITHUB_API_URL:-https://api.github.com}"
HEAD_SHA="$(git rev-parse HEAD)"
HEAD_SHORT="$(git rev-parse --short HEAD)"

bench() {
  go test -run='^$' -bench="$BENCH_PATTERN" -benchmem -count="$COUNT" "${BENCH_PKGS[@]}"
}

echo "=== Benchmarking head ($HEAD_SHORT) ==="
bench | tee head.txt

git fetch --depth=50 origin "$GITHUB_BASE_REF"
BASE_SHORT="$(git rev-parse --short "origin/$GITHUB_BASE_REF")"
git checkout --detach "origin/$GITHUB_BASE_REF"
echo "=== Benchmarking base ($BASE_SHORT) ==="
BASE_OK=1
if ! bench | tee base.txt; then
  echo "base bench run failed; diff will be skipped" >&2
  BASE_OK=0
fi
git checkout --detach "$HEAD_SHA"

go install golang.org/x/perf/cmd/benchstat@latest
BENCHSTAT="$(go env GOPATH)/bin/benchstat"
if [[ "$BASE_OK" == "1" ]]; then
  SUMMARY="$("$BENCHSTAT" base.txt head.txt)"
  HEADER="Benchmark diff: \`$BASE_SHORT\` → \`$HEAD_SHORT\` ($COUNT runs)"
else
  SUMMARY="$("$BENCHSTAT" head.txt)"
  HEADER="Base benchmarks did not run on \`$BASE_SHORT\`; head-only results for \`$HEAD_SHORT\`"
fi
echo "$SUMMARY"
BODY=$(printf "%s\n### %s\n\n\`\`\`\n%s\n\`\`\`\n" "$MARKER" "$HEADER" "$SUMMARY")
if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  echo "$BODY" >> "$GITHUB_STEP_SUMMARY"
fi

if [[ -z "${GITHUB_TOKEN:-}" || -z "$PR_NUMBER" ]]; then
  echo "GITHUB_TOKEN or the number of the pull request is not set, the comment is skipped" >&2
  exit 0
fi
api() {
  curl --fail-with-body -sS -H "Authorization: Bearer $GITHUB_TOKEN" -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" "$@"
}
PAYLOAD=$(jq -n --arg body "$BODY" '{body: $body}')
# the comment of the previous push is updated
COMMENT_ID="$(api "$API_URL/repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments?per_page=100" |
  jq -r --arg marker "$MARKER" '[.[] | select(.body | startswith($marker))][0].id // empty')"
if [[ -n "$COMMENT_ID" ]]; then
  api -X PATCH --data "$PAYLOAD" "$API_URL/repos/$GITHUB_REPOSITORY/issues/comments/$COMMENT_ID" >/dev/null
else
  api -X POST --data "$PAYLOAD" "$API_URL/repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments" >/dev/null
fi
echo "Posted benchmark diff to PR #$PR_NUMBER"
