#!/usr/bin/env bash
# Run `go fmt` on all the source files.

git ls-files | grep --line-buffered '.*\.go$' | xargs realpath | xargs dirname | xargs go fmt
