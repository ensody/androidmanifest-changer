#!/usr/bin/env bash
set -euxo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

source "$DIR/utils.sh"

mkdir -p ~/bin

# Install goreleaser
download_tgz https://github.com/goreleaser/goreleaser/releases/download/v0.182.1/goreleaser_Linux_x86_64.tar.gz \
  bb0b3a96bb38ba86fb3f363d303ce6079c04ada2797a892bed2e2a61ad41daf2 ~/bin \
  goreleaser
