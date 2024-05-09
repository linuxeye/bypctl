package cmder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/scottames/cmder/pkg/log"
)

var dryRun bool

// DryRun will ensure that the commands executed will log, but not be executed
// in the scope of the package
func DryRun() {
	dryRun = true
}

// New returns a new os/exec builder command implemented using the Cmder interface
func New(command ...string) Cmder {
	c := &cmd{strings: command}
	c.ctx = context.Background()
	c.env = os.Environ()
	c.stdin = os.Stdin
	c.stdout = os.Stdout
	c.stderr = os.Stderr

	return c
}

// cmd implements the Cmder interface
type cmd struct {
	cmd       *exec.Cmd
	complete  bool
	ctx       context.Context
	dryRun    bool
	dryRunKey string
	dir       string
	end       time.Time
	env       []string
	exitCode  int
	failed    bool
	logger    log.Logger
	process   *os.Process
	silent    bool
	start     time.Time
	stderr    io.Writer
	stdin     io.Reader
	stdout    io.Writer
	strings   []string
}

func (c *cmd) Args(args ...string) Cmder {
	if len(args) > 0 {
		c.strings = append(c.strings, args...)
	}

	return c
}

func (c *cmd) Clone() Cmder {
	clone := *c
	return &clone
}

func (c *cmd) CombinedOutput() ([]byte, error) {
	var b bytes.Buffer
	c.stdout = &b
	c.stderr = &b
	err := c.Run()

	return b.Bytes(), err
}

func (c *cmd) Complete() bool {
	return c.complete
}

func (c *cmd) Ctx(ctx context.Context) Cmder {
	c.ctx = ctx
	return c
}

func (c *cmd) Dir(dir string) Cmder {
	c.dir = dir
	return c
}

func (c *cmd) DryRun(keys ...string) Cmder {
	c.dryRunKey = strings.Join(keys, " ")
	c.dryRun = true

	return c
}

func (c *cmd) Duration() time.Duration {
	return c.end.Sub(c.start)
}

func (c *cmd) Env(env ...string) Cmder {
	c.env = append(c.env, env...)
	return c
}

func (c *cmd) ExitCode() int {
	return c.exitCode
}

func (c *cmd) In(input ...byte) Cmder {
	if input != nil {
		c.stdin = bytes.NewReader(input)
	} else {
		c.stdin = os.Stdin
	}

	return c
}

func (c *cmd) Kill() error {
	if c.isDryRun() {
		c.logCmdDryRun(log.LoggerKillKey)
		return nil
	}

	c.failed = true
	c.exitCode = -1
	log.LoggerKey = log.LoggerKillKey

	c.LogCmd()

	return c.cmd.Process.Kill()
}

func (c *cmd) LogCmd() {
	if c.silent {
		return
	}

	if c.logger == nil {
		c.logger = getLogger()
	}

	msg := fmt.Sprintf("%v", c.strings)
	if c.dir != "" {
		msg += fmt.Sprintf(string(log.LoggerColor)+" in"+string(log.LoggerClear)+" %s", c.dir)
	}

	c.logger.Log(msg)
}

func (c *cmd) Logger(l log.Logger) Cmder {
	c.logger = l
	return c
}

func (c *cmd) Out(stdout io.Writer, stderr ...io.Writer) Cmder {
	if len(stderr) > 0 {
		c.stderr = stderr[0]
	}

	c.stdout = stdout

	return c
}

func (c *cmd) Output() ([]byte, error) {
	c.Silent()

	if !c.initAndContinue(log.LoggerOutputKey) {
		return nil, nil
	}

	c.clearStdOutStdErr()
	output, err := c.cmd.Output()

	return output, c.endState(err)
}

func (c *cmd) Process() *os.Process {
	return c.process
}

func (c *cmd) Pid() *int {
	if c.process == nil {
		return nil
	}

	return &c.process.Pid
}

func (c *cmd) Run(w ...io.Writer) error {
	if !c.initAndContinue(log.LoggerRunKey, w...) {
		return nil
	}

	err := c.cmd.Run()

	return c.endState(err)
}

func (c *cmd) RunFn(w ...io.Writer) func(args ...string) error {
	return func(args ...string) error {
		return c.Clone().Args(args...).Run(w...)
	}
}

func (c *cmd) RunFnCmd(w ...io.Writer) func(args ...string) (Cmder, error) {
	return func(args ...string) (Cmder, error) {
		clone := c.Clone().Args(args...)
		err := clone.Run(w...)

		return clone, err
	}
}

func (c *cmd) Silent() Cmder {
	c.silent = true
	return c
}

func (c *cmd) Start(w ...io.Writer) error {
	log.LoggerKey = log.LoggerStartKey

	if !c.initAndContinue(log.LoggerStartKey, w...) {
		return nil
	}

	err := c.cmd.Start()
	if err != nil {
		return err
	}

	c.process = c.cmd.Process

	return nil
}

func (c *cmd) StartFn(w ...io.Writer) func(args ...string) error {
	return func(args ...string) error {
		return c.Clone().Args(args...).Start(w...)
	}
}

func (c *cmd) StartFnCmd(w ...io.Writer) func(args ...string) (Cmder, error) {
	return func(args ...string) (Cmder, error) {
		clone := c.Clone().Args(args...)
		err := clone.Start(w...)

		return clone, err
	}
}

func (c *cmd) String() string {
	return c.buildExec().String()
}

func (c *cmd) Wait() error {
	log.LoggerKey = log.LoggerWaitKey

	if c.isDryRun() {
		c.logCmdDryRun(log.LoggerKey)
		return nil
	}

	c.LogCmd()

	if c.process == nil {
		return fmt.Errorf("process expected to be started. found nil process for Wait")
	}

	err := c.cmd.Wait()
	if err != nil {
		c.failed = true
		return err
	}

	c.end = time.Now()
	c.complete = true

	return nil
}

// buildExec builds the exec.Cmd for the given cmd
func (c *cmd) buildExec(w ...io.Writer) *exec.Cmd {
	command := exec.CommandContext(c.ctx, c.strings[0], c.strings[1:]...) //nolint:gosec // written as intended
	c.cmd = command

	switch lw := len(w); {
	case lw == 1:
		c.cmd.Stdout = w[0]
		c.cmd.Stderr = w[0]
	case lw > 1:
		c.cmd.Stdout = w[0]
		c.cmd.Stderr = w[1]
	default:
		c.cmd.Stdout = c.stdout
		c.cmd.Stderr = c.stderr
	}

	c.cmd.Dir = c.dir
	c.cmd.Env = c.env
	c.cmd.Stdin = c.stdin

	return c.cmd
}

// clearStdOutStdErr will set the cmd stdout and stderr to nil
func (c *cmd) clearStdOutStdErr() {
	c.stdout = nil
	c.stderr = nil

	if c.cmd != nil {
		c.cmd.Stdout = nil
		c.cmd.Stderr = nil
	}
}

// endState sets the end state of the command after it's completion
func (c *cmd) endState(err error) error {
	c.end = time.Now()
	c.exitStatus(err)
	c.complete = true

	return err
}

// initAndContinue initializes the cmd and returns a bool value whether it should continue
// or return based on the dryrun value
func (c *cmd) initAndContinue(command string, w ...io.Writer) bool {
	c.buildExec(w...)

	if c.isDryRun() {
		c.logCmdDryRun(command)
		return false
	}

	if !c.silent {
		c.LogCmd()
	}

	c.start = time.Now()

	return true
}

// isDryRun returns whether dryRun is set in the scope of the current cmd or globally
func (c *cmd) isDryRun() bool {
	return dryRun || c.dryRun
}

// exitStatus returns the exit status for the given cmd's error
// if nil found, exit code will always be 0
func (c *cmd) exitStatus(err error) {
	if err == nil {
		c.exitCode = 0
		return
	}

	if e, ok := err.(exitStatus); ok {
		c.exitCode = e.ExitStatus()
	}

	if e, ok := err.(*exec.ExitError); ok {
		if ex, ok := e.Sys().(exitStatus); ok {
			c.exitCode = ex.ExitStatus()
		}
	}

	c.exitCode = 1
}

// exitStatus an interface implementing the ExitStatus method
type exitStatus interface {
	ExitStatus() int
}

// logCmdDryRun logs the given cmd in the context of DryRun
func (c *cmd) logCmdDryRun(s string) {
	// save defaults
	defKey := log.LoggerKey
	defColor := log.LoggerColor

	// override defaults
	log.LoggerColor = log.LoggerDryRunColor
	origSilent := c.unsetSilent()

	if c.dryRunKey != "" {
		log.LoggerKey = c.dryRunKey
	} else {
		log.LoggerKey = log.LoggerDryRunKey + " " + s
		log.LoggerCols = log.LoggerDryRunCols
	}

	// log dry
	c.LogCmd()

	// reset defaults
	log.LoggerKey = defKey
	log.LoggerColor = defColor
	c.silent = origSilent
}

// unsetSilent sets silent to false and returns the original value of c.silent
func (c *cmd) unsetSilent() bool {
	orig := c.silent
	c.silent = false

	return orig
}
