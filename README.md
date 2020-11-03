# circleci-config-merge

Generate .circleci/config.yml by merging multiple files

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

## How to use

```
$ circleci-config-merge merge <file> [<file> ...] > .circleci/config.yml
```

## How to test in CI

In CI, we should test whether `.circleci/config.yml` and the result of `circleci-config-merge merge` is equal as YAML.
`circleci-config-merge` doesn't provide the feature to compare YAML, so please use the other tool like [yamldiff](https://github.com/sahilm/yamldiff).

## Split File Format

The split file format is same as [.circleci/config.yml](https://circleci.com/docs/2.0/configuration-reference/).

## Merge Rule

Coming soon.

## Example

Coming soon.

## LICENSE

[MIT](LICENSE)
