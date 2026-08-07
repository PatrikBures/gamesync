#!/usr/bin/env bash

count=$(find . -type f -name '*.go' \
    -not -path './internal/ogen/*' \
    -not -path './internal/server/dbm/*' \
    -not -path './internal/server/permissions/perm_string.go' \
    -exec cat {} \; |
    grep -cv '^$'
)

echo "lines written in go excluding generated and empty lines:"
echo "$count"
