#!/usr/bin/env bash

set -euo pipefail

ROOT="$(dirname "$(cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd)")"

FROM=golang.org/x/tools
PKG=github.com/charlievieth/xtools
TO="$PKG"

git fetch --all
git checkout origin/golang internal

for file in ./internal/*; do
    mv "$file" "./$(basename "$file")"
done
