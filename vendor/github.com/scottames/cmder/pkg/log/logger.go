package log

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// Logger is a generic logging interface
// see also: github.com/go-log/log
//
// By default the built-in logger will be used. See Additional Logger* variables
// for configuration options.
//
// Color is disabled by default, but can be enabled by setting either
// MAGEFILE_ENABLE_COLOR or CMDER_ENABLE_COLOR environment variables to true.
type Logger interface {
	// Log inserts a log entry.  Arguments may be handled in the manner
	// of fmt.Print, but the underlying logger may also decide to handle
	// them differently.
	Log(v ...interface{})

	// Logf inserts a log entry.  Arguments are handled in the manner of
	// fmt.Printf.
	Logf(format string, v ...interface{})
}

// Color a string alias for logging colors
type Color string

var (
	// LoggerKey the key used to when logging the command
	// if the default logger is used and Silent is not specified
	// represents the action of the command
	LoggerKey = LoggerRunKey

	// LoggerCols specifies the right justified columns the LoggerKey will be padded
	LoggerCols = "5"

	// LoggerColor the default color used to print the LoggerKey
	// if the default logger is used and Silent is not specified
	LoggerColor = LoggerTeal

	// LoggerDryRunColor the color used to print the LoggerDryRunKey when DryRun
	// is invoked if the default logger is used and Silent is not specified
	LoggerDryRunColor = LoggerYellow

	// LoggerDryRunCols specifies the right justified columns the LoggerKey will be padded when
	// DryRun is invoked
	LoggerDryRunCols = "10"

	// LoggerDryRunKey the key used to represent the action of the command when DryRun is invoked
	LoggerDryRunKey = "dry"

	// LoggerKillKey the key used to represent the action of the command when Kill is invoked
	LoggerKillKey = "kill"

	// LoggerOutputKey the key used to represent the action of the command when Output is invoked
	LoggerOutputKey = "output"

	// LoggerRunKey the key used to represent the action of the command when Run is invoked
	LoggerRunKey = "run"

	// LoggerStartKey the key used to represent the action of the command when Start is invoked
	LoggerStartKey = "start"

	// LoggerWaitKey the key used to represent the action of the command when Start is invoked
	LoggerWaitKey = "wait"

	// LoggerBlack the default color code representing the color Black in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerBlack Color = "\033[1;30m"

	// LoggerRed the default color code representing the color Red in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerRed Color = "\033[1;31m"

	// LoggerGreen the default color code representing the color Green in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerGreen Color = "\033[1;32m"

	// LoggerYellow the default color code representing the color Yellow in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerYellow Color = "\033[1;33m"

	// LoggerPurple the default color code representing the color Purple in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerPurple Color = "\033[1;34m"

	// LoggerMagenta the default color code representing the color Magenta in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerMagenta Color = "\033[1;35m"

	// LoggerTeal the default color code representing the color Teal in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerTeal Color = "\033[1;36m"

	// LoggerWhite the default color code representing the color White in human readable format
	// only used when set to either LoggerColor or LoggerDryRunColor
	LoggerWhite Color = "\033[1;37m"

	// LoggerDarkGrey the default color code representing the color DarkGrey in human readable
	// format only used when set to either LoggerColor or LoggerDryRunColor
	LoggerDarkGrey Color = "\033[90m"

	// LoggerClear the default color code used to clear colors or return the print statement to the default
	LoggerClear Color = "\033[0m"

	// LoggerTimeFormat format to use when logging the timestamp - when LoggerTimeStampEnabled
	LoggerTimeFormat = time.RFC3339

	// LoggerTimeStampEnabled if set to true logger will include a timestamp
	LoggerTimeStampEnabled = true
)

func init() {
	var enableColor bool

	if envBool("MAGEFILE_ENABLE_COLOR") || envBool("CMDER_ENABLE_COLOR") {
		enableColor = true
	}

	if !enableColor {
		LoggerColor = ""
		LoggerDryRunColor = ""

		LoggerBlack = ""
		LoggerRed = ""
		LoggerGreen = ""
		LoggerYellow = ""
		LoggerPurple = ""
		LoggerMagenta = ""
		LoggerTeal = ""
		LoggerWhite = ""

		LoggerDarkGrey = ""
		LoggerClear = ""
	}
}

// New returns a new logger instance which implements the Logger interface
func New() *logger { //nolint:revive // the intention is to leverage the methods and interface
	return &logger{key: LoggerKey, color: LoggerColor}
}

// logger implements the Logger interface and adds additional functionality
type logger struct {
	color       Color
	key         string
	noTimestamp bool
}

// Key sets the logger key for the given logger instance
func (l *logger) Key(k string) *logger {
	l.key = k
	return l
}

// Color sets the logger key for the given logger instance
func (l *logger) Color(lc Color) *logger {
	l.color = lc
	return l
}

// WithoutTimestamp sets the current logger instance to omit the timestamp when logging
func (l *logger) WithoutTimestamp() *logger {
	l.noTimestamp = true
	return l
}

// Logf implements the Logger interface
func (l logger) Logf(format string, v ...interface{}) {
	s := l.prependStr()
	timestamp := l.timestamp(string(LoggerDarkGrey))
	fmt.Printf(s+format+timestamp+"\n", v...)
}

// splitArgsToNewLine returns a new string formatted as a shell command as if being executed
// on multiple lines - padding between the [] of the command (slice) to be printed
func splitArgsToNewLine(s string) string {
	var result string

	const (
		bBracket = "["
		eBracket = "]"
	)

	newArgs := []string{}
	args := strings.Split(s, " ")

	for _, w := range args {
		if strings.HasPrefix(w, bBracket) {
			w = strings.Replace(w, bBracket, bBracket+"\n\n ", 1)
		}

		if strings.HasSuffix(w, eBracket) {
			i := strings.LastIndex(w, eBracket)
			w = w[:i] + strings.Replace(w[i:], eBracket, "\n\n"+eBracket, 1)
		}

		if strings.Contains(w, eBracket+string(LoggerColor)) {
			w = strings.Replace(w, eBracket+string(LoggerColor), "\n\n"+eBracket+string(LoggerColor), 1)
		}

		if strings.HasPrefix(w, "-") {
			newArgs = append(newArgs, "\\\n   "+w)
		} else {
			newArgs = append(newArgs, w)
		}
	}

	result = strings.Join(newArgs, " ")

	return result
}

// stringLen returns the length of a given string
func stringLen(s string) int {
	return len([]rune(s))
}

// Log implements the Logger interface
func (l logger) Log(v ...interface{}) {
	s := l.prependStr()
	timestamp := l.timestamp(string(LoggerDarkGrey))
	msg := fmt.Sprintf("%v", v...)

	termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err == nil {
		absLen := stringLen(msg) + stringLen(s) + stringLen(timestamp)
		if absLen > termWidth {
			msg = splitArgsToNewLine(msg)
		}
	}

	fmt.Printf(s+"%s"+timestamp+"\n", msg)
}

func (l logger) prependStr() string {
	key := l.getKey()
	colonColor := string(LoggerDarkGrey)

	return fmt.Sprintf(
		string(l.getColor())+
			"%"+
			LoggerCols+
			"s"+
			colonColor+
			" : "+
			string(LoggerClear),
		key,
	)
}

func (l logger) timestamp(color string) string {
	if l.noTimestamp || !LoggerTimeStampEnabled {
		return ""
	}

	now := time.Now().Format(LoggerTimeFormat)

	return color + " : " + now + string(LoggerClear)
}

func (l logger) getColor() Color {
	if l.color != LoggerColor {
		return l.color
	}

	return LoggerColor
}

func (l logger) getKey() string {
	if l.key != LoggerKey {
		return l.key
	}

	return LoggerKey
}

// envBool returns true if the given env variable is set to a truth
// as defined by strconv.ParseBool
func envBool(env string) bool {
	b, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		return false
	}

	return b
}
