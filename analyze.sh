#!/bin/bash

for dir in $(find . -type d); do
    if ls "$dir"/*.go >/dev/null 2>&1; then
        echo "Analyzing package: $dir"
        ./bin/staticlint "$dir"
    fi
done