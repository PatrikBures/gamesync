#!/usr/bin/env sh

set -e

# 512000 = 500KiB
max=512000

human() {
    # takes in bytes
    if [ "$1" -lt 1024 ]; then
        echo "$1"
    else
        echo "$(($1 / 1024)) KiB"
    fi
}

merge_from_branch=$1
merge_to_branch=$2

if [ "$#" -ne 2 ] || [ -z "$merge_from_branch" ] || [ -z "$merge_to_branch" ]; then
    echo "usage: $0 MERGE_FROM_BRANCH MERGE_TO_BRANCH"
    exit 2
fi

echo "from: $merge_from_branch"
echo "to:   $merge_to_branch"
echo

size_diff_bytes=$(git rev-list --disk-usage --objects "$merge_to_branch..$merge_from_branch")

if [ "$size_diff_bytes" -gt "$max" ]; then
    echo "--- BRANCH CHANGES ARE TOO BIG!!! ---"
    echo "branch diff size is $(human "$size_diff_bytes")."
    echo "which is $(human $((size_diff_bytes - max))) bytes too many."
    echo "you should probably ask the git-master for help."
    exit 1
else
    echo "size diff in the clear."
    echo "size diff is $(human "$size_diff_bytes")."
    echo "you can use $(human $((max - size_diff_bytes))) more."
fi

echo
total_to=$(git rev-list --disk-usage --objects "$merge_to_branch")
total_from=$(git rev-list --disk-usage --objects "$merge_from_branch")
echo "total size of $merge_to_branch is $(human "$total_to")"
echo "total size of $merge_from_branch is $(human "$total_from")"
