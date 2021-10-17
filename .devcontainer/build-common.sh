#!/usr/bin/env bash
set -euxo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

source "$DIR/utils.sh"

mkdir -p ~/bin

# Install goreleaser
download_tgz https://github.com/goreleaser/goreleaser/releases/download/v0.174.2/goreleaser_Linux_x86_64.tar.gz \
  38155642fb10a75205f20e390474f3bad9fbf61f2614500b02b179d05907348e ~/bin \
  goreleaser
