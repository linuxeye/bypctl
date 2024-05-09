# Cmder

[![Lint & Test](https://github.com/scottames/cmder/actions/workflows/pull_request.yml/badge.svg)](https://github.com/scottames/cmder/actions/workflows/pull_request.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/scottames/cmder.svg)](https://pkg.go.dev/github.com/scottames/cmder)
[![Go Report Card](https://goreportcard.com/badge/github.com/scottames/cmder)](https://goreportcard.com/report/github.com/scottames/cmder)

Cmder implements builder functionality wrapping os/exec for easily constructing shell commands. The
initial intended use for this package was simple scripting with [Mage](https://github.com/magefile/mage), but there is nothing that
should stop it from being used in other implementations.

Heavily inspired by [Mage](https://github.com/magefile/mage)'s own [sh](https://github.com/magefile/mage/tree/master/sh) library and many more great projects who have come before.

## Usage

### Examples

Examples can be found in the [`_examples/magefile.go`](_examples/magefile.go)
[Magefile](https://magefile.org/).

Run the following to execute any of the example targets from the above Magefile.

```shell
cd _examples
mage -l
```

Additionally many of the tests can provide some example usage.

### Logging

Cmder logs all commands being run, using the `Logger` method, which implements the [`Logger`](https://github.com/scottames/cmder/blob/master/pkg/log/logger.go#L10-L27) interface:

```golang
type Logger interface {
  // Log inserts a log entry. Arguments may be handled in the manner
  // of fmt.Print, but the underlying logger may also decide to handle
  // them differently.
  Log(v ...interface{})

  // Logf inserts a log entry. Arguments are handled in the manner of
  // fmt.Printf.
  Logf(format string, v ...interface{})
}
```

By default (if none specified with the `Cmder.Logger()` method) the built-in [logger](pkg/log/logger.go) will be used. See Additional `log.Logger*` variables for configuration
options.

Color is disabled by default, but can be enabled by setting either `MAGEFILE_ENABLE_COLOR` or
`CMDER_ENABLE_COLOR` environment variables to true.

The default logger will check the terminal width and if the command to be printed is wider than the terminal width, it will be broken up into multiple lines, similar to a shell command represented on multiple lines.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Contributors should follow the [Go Community Code of Conduct
](https://golang.org/conduct).

### Tests

Requires:

- [mage](https://magefile.org/)
- [golangci-lint](https://golangci-lint.run/)

Run:

```shell
mage check
```
