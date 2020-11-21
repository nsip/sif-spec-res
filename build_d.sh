#!/bin/bash
set -e

CGO_ENABLED=0 go run ./txt2toml/main.go
CGO_ENABLED=0 go run ./toml2json/config.go ./toml2json/main.go
