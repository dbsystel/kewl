# Contribution Guidelines

## General

Before you push a PR, please create an [issue](../../issues).

That way we can manage and track the changes being made to the project.

Make sure your code has been linted successfully, and none of the tests fail.

We would kindly ask you to not decrease the code coverage, so please write/adapt tests accordingly.

## Communication

We do not tolerate violent/racist/sexist or any other behavior aiming to harm anyone, so please respect each other as human beings. Thank you.

## Linting

To lint the project, please proceed as follows:

- install [golangci-lint](https://golangci-lint.run/usage/install/)
- run the linter using `golangci-lint run ./pkg/... --fix`
- fix the remaining errors (if any)

## Test coverage

To test with coverage, please proceed as follows:

- install [go-acc](https://github.com/ory/go-acc) using `go get github.com/ory/go-acc`
- run `go-acc ./pkg/...`

## Releases

Releases are being made by merging the ``develop`` or issue-branch to `master` and tag the version correctly.
