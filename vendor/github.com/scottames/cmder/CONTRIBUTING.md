# Contributing

This may be a small personal project, but if you find it useful, contributions are more than
welcome.

Prior to any commit/pull request please run testing & linting. This can be accomplished by using the
`mage check` target.

## Issues

Please always create an issue before sending a PR unless it's an obvious typo
or other trivial change.

## Dependency Management

Cmder is intended to only use the standard libary as it is intended to be a small wrapper around
os/exec.

## Testing

Please write tests for any new features. Testing can be
accomplished by using the `mage test` target.

## Linting

[Golangci-lint](https://golangci-lint.run/) is used to lint this project. Linting can be
accomplished by using the `mage lint` target.
