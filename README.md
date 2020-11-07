# circleci-config-merge

[![Build Status](https://github.com/suzuki-shunsuke/circleci-config-merge/workflows/CI/badge.svg)](https://github.com/suzuki-shunsuke/circleci-config-merge/actions)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b34ffd9a1198b2952d46/test_coverage)](https://codeclimate.com/github/suzuki-shunsuke/circleci-config-merge/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/circleci-config-merge)](https://goreportcard.com/report/github.com/suzuki-shunsuke/circleci-config-merge)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/circleci-config-merge.svg)](https://github.com/suzuki-shunsuke/circleci-config-merge)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/circleci-config-merge/master/LICENSE)

Generate .circleci/config.yml by merging multiple files

## Blog

https://dev.to/suzukishunsuke/splitting-circleci-config-yml-10gk

## Motivation

Our monivation is to split a huge .circleci/config.yml per service.
We have a [monorepo](https://en.wikipedia.org/wiki/Monorepo) where many services are managed.
`.circleci/config.yml` of this repository has over 6000 lines, and it's hard to maintain the file.

## Why don't we use `circleci config pack`?

The directory structure and naming rule don't match our needs.
And we want to merge the list of workflow's jobs.

## Install

Download from [GitHub Releases](https://github.com/suzuki-shunsuke/circleci-config-merge/releases)

```
$ circleci-config-merge --version
circleci-config-merge version 0.1.0
```

## Example

Please see [examples](examples) and [example-cirleci-config-merge](https://github.com/suzuki-shunsuke/example-circleci-config-merge).

## How to use

```
$ circleci-config-merge merge <file> [<file> ...] > .circleci/config.yml
```

## How to test in CI

In CI, we should test whether `.circleci/config.yml` and the result of `circleci-config-merge merge` is equal as YAML.
`circleci-config-merge` doesn't provide the feature to compare YAML, so please use the other tool like [dyff](https://github.com/homeport/dyff).
Please see the example [suzuki-shunsuke/example-circleci-config-merge](https://github.com/suzuki-shunsuke/example-circleci-config-merge) as a reference to split .circleci/config.yml and setup CI.

## Split File Format

The split file format is same as [.circleci/config.yml](https://circleci.com/docs/2.0/configuration-reference/).

## Merge Rule

Coming soon.

## LICENSE

[MIT](LICENSE)
