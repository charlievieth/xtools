#!/usr/bin/env bash

set -euo pipefail

ROOT="$(dirname "$(cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd)")"

FROM=golang.org/x/tools
PKG=github.com/charlievieth/xtools
TO="$PKG"

GREP_FLAGS=(
    --recursive
    --exclude-dir 'scripts'
    --exclude-dir '.git'
    --files-with-matches
    --null
    --fixed-strings
)

escape_import_path() {
    sed -e 's/\./\\\./g' -e 's/\//\\\//g' <<<"$1"
}

fix_import_paths() {
    local from="$1"
    local to="$2"
    local dir="${3:-$ROOT}"

    local replace
    replace="s/$(escape_import_path "$from")/$(escape_import_path "$to")/g"

    if ! grep "${GREP_FLAGS[@]}" "$from" "$dir" | xargs -0 -- sed -i "$replace"; then
        echo "error: fix_import_paths" >&2
        return 1
    fi
    return 0
}

git fetch --all

LATEST_SHA="$(git rev-parse golang/master)"
BRANCH="golang-$(head -c 8 <<<"$LATEST_SHA")"

# if git branch | grep -F "$BRANCH"; then
#     git branch -D "$BRANCH"
# fi
# git checkout -b "$BRANCH"

git checkout golang/master internal
git restore --staged ./internal

fix_import_paths golang.org/x/tools/internal github.com/charlievieth/xtools ./internal

# PACKAGES=()

for file in ./internal/*; do
    dest="./$(basename "$file")"
    mv "$file" "$dest"
    git add "$dest"
    # if [[ -d $dest ]]; then
    #     PACKAGES+=("$dest")
    # fi
done
rm -r ./internal

# for pkg in "${PACKAGES[@]}"; do
#     fix_import_paths "golang.org/x/tools/internal/$pkg" "github.com/charlievieth/xtools/$pkg"
# done

git checkout golang/master gopls
git restore --staged ./gopls
rm ./gopls/go.{mod,sum}

fix_import_paths golang.org/x/tools/internal github.com/charlievieth/xtools ./gopls
fix_import_paths golang.org/x/tools/gopls/internal github.com/charlievieth/xtools/gopls ./gopls
fix_import_paths golang.org/x/tools/gopls github.com/charlievieth/xtools/gopls ./gopls

# Fix main.go
fix_import_paths 'package main // import "golang.org/x/tools/gopls"' 'package main' ./gopls
MAIN="$(cat ./gopls/main.go)"
printf '// +build never\n\n%s' "$MAIN" > ./gopls/main.go

for file in ./gopls/internal/*; do
    dest="./gopls/$(basename "$file")"
    mv "$file" "$dest"
    git add "$dest"
done
rm -r ./gopls/internal

git add ./gopls

# GO111MODULE=on go mod tidy
# git add go.mod go.sum

git commit -m "update golang.org/x/tools to $(head -c 8 <<<"$LATEST_SHA")"
