package cmder

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/scottames/cmder/pkg/log"
)

// Cmder the os/exec shell command builder interface
type Cmder interface {
	// Args appends additional arguments to the given command
	Args(...string) Cmder

	// Ctx can be used to pass context to the underlying exec command. The Run method will call
	// exec.CommandContext with the given context
	Ctx(context.Context) Cmder

	// CombinedOutput runs the command and returns its combined
	// standard output and standard error.
	CombinedOutput() ([]byte, error)

	// Complete returns a boolean value as to whether the command has complete it's execution
	Complete() bool

	// Clone returns a new Cmder copied from the original
	Clone() Cmder

	// Dir specifies the working directory of the command.
	//
	// If Dir is not set, Run runs the command in the
	// calling process's current directory.
	Dir(string) Cmder

	// DryRun sets the current Cmder instance to run in DryRun mode (log and skip execution)
	//
	// Option to pass one or more strings to replace the dryrun action key
	// defaults to LoggerDryRunKey
	DryRun(...string) Cmder

	// Duration returns the duration in time.Duration the command took to run
	// if it has completed.
	Duration() time.Duration

	// Env appends the given strings to the environment of the process
	//
	// Each entry is of the form "key=value".
	//
	// The new process uses the calling process's environment in addition
	// to the keys provided by Env.
	//
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	Env(...string) Cmder

	// ExitCode returns the exit code of the command. If Run has not been invoked
	// zero will always be returned.
	ExitCode() int

	// In connects the new process' stdin to the current process's stdin if no input provided.
	// If input provided the input will be passed to the new process' stdin.
	In(...byte) Cmder

	// Out connects the new process' stdout and optionally stderr to the given io.Writers
	// Useful for writing to buffer or files
	Out(stdout io.Writer, stderr ...io.Writer) Cmder

	// Kill invokes the os.exec Kill method on the command
	//
	// Kill causes the Process to exit immediately.
	// Kill does not wait until the Process has actually exited.
	// This only kills the Process itself, not any other processes it may have started.
	Kill() error

	// LogCmder will print the command that is to be executed.
	// Included in Run if Silent unset
	// See also String, LogCmd
	LogCmd()

	// Allows setting an external logger.
	// See the Logger interface for more details.
	Logger(log.Logger) Cmder

	// Output invokes the os.exec Output method on the command
	//
	// Output runs the command and returns its stdout.
	// Any returned error will usually be of type *exec.ExitError.
	Output() ([]byte, error)

	// Pid returns the process id of the exited process or nil if the process has yet to exit.
	// See also
	// - Process
	// - https://pkg.go.dev/os#Process
	Pid() *int

	// Process is the underlying process, once started.
	// See also: https://pkg.go.dev/os#Process
	Process() *os.Process

	// Run invokes the os.exec Run method on the command
	//
	// Run calls exec.Run starting the specified command and waits for it to complete.
	// If Silent not set, the command will be printed prior to execution.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	// The returned error is nil if the command runs, has no problems copying stdin, stdout, and stderr, and exits with a zero exit status.
	//
	// If the command starts but does not complete successfully, the error is of type *ExitError. Other error types may be returned for other situations.
	//
	// If the calling goroutine has locked the operating system thread with runtime.LockOSThread and modified any inheritable OS-level thread state (for example, Linux or Plan 9 name spaces), the new process will inherit the caller's thread state.
	Run(...io.Writer) error

	// RunFn returns a function to call Run with the given command
	// optionally appending any new args and returning an error.
	// When the function is invoked a new underlying command is created
	// in order to retain the state of the original.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	// Can be useful in defining a command once and calling it
	// multiple times or ways.
	//
	// See also: RunFnWithCmd
	RunFn(...io.Writer) func(args ...string) error

	// RunFnCmd returns a function to call Run with the given command
	// optionally appending any new args and returning the new command
	// and a possible error.
	// When the function is invoked a new underlying command is created
	// in order to retain the state of the original.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	// Can be useful in defining a command once and calling it
	// multiple times or ways.
	//
	// See also: RunFn
	RunFnCmd(...io.Writer) func(args ...string) (Cmder, error)

	// Silent will set Run to not print the command prior to execution
	Silent() Cmder

	// Start invokes the os.exec Start method on the command
	//
	// Start starts the specified command but does not wait for it to complete.
	//
	// If Start returns successfully, the Process method will be set.
	//
	// The Wait method will return the exit code and release associated resources
	// once the command exits.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	Start(...io.Writer) error

	// StartFn returns a function to call Run with the given command
	// optionally appending any new args and returning an error.
	// When the function is invoked a new underlying command is created
	// in order to retain the state of the original.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	// Can be useful in defining a command once and calling it
	// multiple times or ways.
	//
	// See also: StartFnCmd, RunFun, RunFnWithCmd
	StartFn(w ...io.Writer) func(args ...string) error

	// StartFnCmd returns a function to call Run with the given command
	// optionally appending any new args and returning the new command
	// and a possible error.
	// When the function is invoked a new underlying command is created
	// in order to retain the state of the original.
	//
	// Optionally one or two io.Writer may be passed to Run where
	// the first being stdout and the second being stderr.
	// If only one specified it will be connected to the processes' stdout and stderr.
	//
	// Can be useful in defining a command once and calling it
	// multiple times or ways.
	//
	// See also: StartFn, RunFnCmd, RunFn
	StartFnCmd(...io.Writer) func(args ...string) (Cmder, error)

	// String returns a human-readable description of the command
	// from exec.Cmder. It is intended only for debugging.
	// In particular, it is not suitable for use as input to a shell.
	// The output of String may vary across Go releases.
	// See also: DryRun
	String() string

	// Wait invokes the os.exec Wait method on the command
	//
	// Wait waits for the command to exit and waits for any copying to
	// stdin or copying from stdout or stderr to complete.
	//
	// The command must have been started by Start.
	//
	// The returned error is nil if the command runs, has no problems
	// copying stdin, stdout, and stderr, and exits with a zero exit
	// status.
	//
	// If the command fails to run or doesn't complete successfully, the
	// error is of type *ExitError. Other error types may be
	// returned for I/O problems.
	//
	// If any of c.Stdin, c.Stdout or c.Stderr are not an *os.File, Wait also waits
	// for the respective I/O loop copying to or from the process to complete.
	//
	// Wait releases any resources associated with the Cmd.
	Wait() error
}
