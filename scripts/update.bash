#!/usr/bin/env bash

set -euo pipefail
# WARN
# set -x

ROOT="$(dirname "$(cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd)")"

# change to root directory
cd "$ROOT"

FROM=golang.org/x/tools
PKG=github.com/charlievieth/xtools
TO="$PKG"


escape_import_path() {
    sed -e 's/\./\\\./g' -e 's/\//\\\//g' <<<"$1"
}

fix_import_paths() {
    local from="$1"
    local to="$2"
    local dir="${3:-$ROOT}"

    local replace
    replace="s/$(escape_import_path "$from")/$(escape_import_path "$to")/g"

    GREP_FLAGS=(
        --recursive
        --exclude-dir 'scripts'
        --exclude-dir '.git'
        --files-with-matches
        --null
        --fixed-strings
    )
    if ! grep "${GREP_FLAGS[@]}" "$from" "$dir" | xargs -0 -- sed -i "$replace"; then
        echo "error: fix_import_paths" >&2
        return 1
    fi
    return 0
}

disable_gopls_main() {
    fix_import_paths 'package main // import "golang.org/x/tools/gopls"' 'package main'
    local src
    src="$(cat ./gopls/main.go)"
    disable='func init() { panic("WRONG GOPLS") }'
    printf '// +build never\n\n%s\n\n%s\n' "${src}" "${disable}" > ./gopls/main.go
    gofmt -w ./gopls/main.go
    git add ./gopls/main.go
}

git fetch --all

LATEST_SHA="$(git rev-parse golang/master)"
BRANCH="golang-$(head -c 8 <<<"$LATEST_SHA")"

if git branch | grep -F "$BRANCH"; then
    git branch -D "$BRANCH" # WARN
fi
git checkout -b "$BRANCH"

git checkout golang/master internal
git restore --staged ./internal

fix_import_paths golang.org/x/tools/internal github.com/charlievieth/xtools ./internal

for file in ./internal/*; do
    dest="./$(basename "$file")"
    if [[ -d "$dest" ]]; then
        git rm -rf "$dest"
    elif [[ -f "$dest" ]]; then
        git rm -f "$dest"
    fi
    mv "$file" "$dest"
    git add "$dest"
done
rm -r ./internal

rm -r ./gopls
git checkout golang/master gopls
git restore --staged ./gopls

# remove gopls mod file
rm ./gopls/go.{mod,sum}

# disable the gopls main (so we don't install it by accident)
disable_gopls_main

fix_import_paths golang.org/x/tools/internal github.com/charlievieth/xtools ./gopls
fix_import_paths golang.org/x/tools/gopls/internal github.com/charlievieth/xtools/gopls ./gopls
fix_import_paths golang.org/x/tools/gopls github.com/charlievieth/xtools/gopls ./gopls

for file in ./gopls/internal/*; do
    dest="./gopls/$(basename "$file")"
    if [[ -e "$dest" ]]; then
        # This should not happen since we first delete the
        # gopls directory
        echo "WARN: overwriting: $dest"
    fi
    rm -rf "$dest"
    mv "$file" "$dest"
    git add "$dest"
done
rm -r ./gopls/internal

git add ./gopls

# pull the latest tools since gopls depends on them
GO111MODULE=on go get golang.org/x/tools@master
GO111MODULE=on go mod tidy
git add go.mod go.sum

# Make sure we can build gopls
(
    cd ./gopls
    go build -tags never
    rm ./gopls
)

git commit -m "update golang.org/x/tools to $(head -c 8 <<<"$LATEST_SHA")

https://github.com/golang/tools/commit/${LATEST_SHA}"
