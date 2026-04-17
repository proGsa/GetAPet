#!/bin/sh
set -eu

TMP_OUT="$(mktemp)"
trap 'rm -f "$TMP_OUT"' EXIT

if go test -tags=integration -json ./tests/unit -count=1 >"$TMP_OUT" 2>&1; then
  TEST_EXIT=0
else
  TEST_EXIT=1
fi

cat "$TMP_OUT"

PASSED_COUNT=$(grep -c "PASS" "$TMP_OUT" || true)
FAILED_COUNT=$(grep -c "FAIL" "$TMP_OUT" || true)

echo "========================================"
echo "Unit tests summary:"
echo "Passed lines: $PASSED_COUNT"
echo "Failed lines: $FAILED_COUNT"
echo "========================================"

exit "$TEST_EXIT"