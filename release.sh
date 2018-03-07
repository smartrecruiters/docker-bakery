#!/usr/bin/env bash
# GITHUB_TOKEN needs to be available for the goreleaser
# export GITHUB_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# https://goreleaser.com/#quick_start

git tag -a v -m "${1} Release"
git push origin v${1}

goreleaser --snapshot --rm-dist