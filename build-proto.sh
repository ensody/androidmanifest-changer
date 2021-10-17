#!/usr/bin/env bash
set -euxo pipefail

protoc --go_out=. *.proto
