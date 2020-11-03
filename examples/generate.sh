#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")"

cat header.txt > .circleci/config.yml
bash merge.sh >> .circleci/config.yml
