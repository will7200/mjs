package hooks

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"runtime"

	"github.com/sirupsen/logrus"
)

type CallerHook struct {
	CallerHookOptions *CallerHookOptions
}

// CREDIT TO magic53
// https://github.com/sirupsen/logrus/pull/544

// NewHook creates a new caller hook with options. If options are nil or unspecified, options.Field defaults to "src"
// and options.Flags defaults to log.Llongfile
func NewHook(options *CallerHookOptions) *CallerHook {
	// Set default caller field to "src"
	if options.Field == "" {
		options.Field = "src"
	}
	// Set default caller flag to Std logger log.Llongfile
	if options.Flags == 0 {
		options.Flags = log.Llongfile
	}
	return &CallerHook{options}
}

// CallerHookOptions stores caller hook options
type CallerHookOptions struct {
	// Field to display caller info in
	Field string
	// Stores the flags
	Flags int
}

// HasFlag returns true if the report caller options contains the specified flag
func (options *CallerHookOptions) HasFlag(flag int) bool {
	return options.Flags&flag != 0
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.CallerHookOptions.Field] = hook.callerInfo() // add 1 for this frame
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CallerHook) callerInfo() string {
	// Follows output of Std logger
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!goSrcRegexp.MatchString(file) || goTestRegexp.MatchString(file)) {
			if hook.CallerHookOptions.HasFlag(log.Lshortfile) && !hook.CallerHookOptions.HasFlag(log.Llongfile) {
				file = path.Base(file)
			}
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	file := "???"
	line := 0
	return fmt.Sprintf("%s:%d", file, line)
}

var goSrcRegexp = regexp.MustCompile(`logrus/.*.go|will7200.*.hooks/.*.go`)
var goTestRegexp = regexp.MustCompile(`.*test.go`)
