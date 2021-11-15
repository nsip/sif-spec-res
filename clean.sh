#!/bin/bash

set -e

rm -rf ./3.4.2/json/ ./3.4.2/toml/
rm -rf ./3.4.3/json/ ./3.4.3/toml/
rm -rf ./3.4.4/json/ ./3.4.4/toml/
rm -rf ./3.4.5/json/ ./3.4.5/toml/
rm -rf ./3.4.6/json/ ./3.4.6/toml/
rm -rf ./3.4.7/json/ ./3.4.7/toml/
rm -rf ./3.4.8/json/ ./3.4.8/toml/
rm -rf ./3.4.9/json/ ./3.4.9/toml/

# rm -rf ./toml2json/config.go  # ---  keep it to avoid 'not found' error

rm -f ./sif-spec-res