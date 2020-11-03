#!/usr/bin/env bash

set -eu
set -o pipefail

cd "$(dirname "$0")"

split_files() {
  find . -path "*/.circleci/*.yml" | grep -v "^\./\.circleci/"
}

split_files | xargs circleci-config-merge merge
