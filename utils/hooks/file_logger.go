package hooks

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
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
	entry.Data[hook.CallerHookOptions.Field] = hook.callerInfo(7 + 1) // add 1 for this frame
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CallerHook) callerInfo(skipFrames int) string {
	// Follows output of Std logger
	_, file, line, ok := runtime.Caller(skipFrames)
	if !ok {
		file = "???"
		line = 0
	} else {
		// check flags
		if hook.CallerHookOptions.HasFlag(log.Lshortfile) && !hook.CallerHookOptions.HasFlag(log.Llongfile) {
			file = path.Base(file)
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
