#!/usr/bin/env bash
set -euxo pipefail

protoc --go_out=. --go-vtproto_out=. --go-vtproto_opt=features=marshal+unmarshal+size *.proto
