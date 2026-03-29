#!/bin/sh
set -eu

TMP_OUT="$(mktemp)"
trap 'rm -f "$TMP_OUT"' EXIT

if go test -tags=integration -json ./tests/integration -count=1 >"$TMP_OUT" 2>&1; then
  TEST_EXIT=0
else
  TEST_EXIT=1
fi

cat "$TMP_OUT"

PASSED_COUNT="$(grep -Ec '"Action":"pass".*"Test":"' "$TMP_OUT" || true)"
FAILED_COUNT="$(grep -Ec '"Action":"fail".*"Test":"' "$TMP_OUT" || true)"

echo "========================================"
echo "Integration tests summary:"
echo "Passed: $PASSED_COUNT"
echo "Failed: $FAILED_COUNT"
echo "========================================"

exit "$TEST_EXIT"