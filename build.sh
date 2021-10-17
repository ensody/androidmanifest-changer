#!/usr/bin/env bash
set -euxo pipefail

goreleaser --snapshot --skip-publish --rm-dist
