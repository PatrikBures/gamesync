#!/usr/bin/env sh

ignored=$(git ls-files -i -c --exclude-standard)
if ! [ "$ignored" = "" ]; then 
    echo "--- FILES THAT ARE IGNORED ARE IN REPO!!! ---"
    echo "files which were found in repo:"
    echo "$ignored"
    exit 1
else
    echo "did not find any files/dirs in the repo which should be ignored"
fi
