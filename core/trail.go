package core

import (
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/colors"
	"log"
	"runtime/debug"
	"strings"
)

// Reporting Levels
const (
	DEBUG     = 0
	WORKING   = 1
	INFO      = 2
	OK        = 3
	WARNING   = 4
	ERROR     = 5
	CRITICAL  = 6
	ALERT     = 7
	EMERGENCY = 8
)

var trailTag = map[int]string{
	DEBUG:     colors.Debug,
	WORKING:   colors.Working,
	INFO:      colors.Info,
	OK:        colors.OK,
	WARNING:   colors.Warning,
	ERROR:     colors.Error,
	CRITICAL:  colors.Critical,
	ALERT:     colors.Alert,
	EMERGENCY: colors.Emergency,
}

var levelMap = map[int]string{
	DEBUG:     "[  DEBUG ]   ",
	WORKING:   "[ WORKING]   ",
	INFO:      "[  INFO  ]   ",
	OK:        "[   OK   ]   ",
	WARNING:   "[ WARNING]   ",
	ERROR:     "[  ERROR ]   ",
	CRITICAL:  "[CRITICAL]   ",
	ALERT:     "[  ALERT ]   ",
	EMERGENCY: "[  EMERG ]   ",
}

// ReportingLevel is the standard reporting level.
var ReportingLevel = DEBUG

// Trail prints to the log
func Trail(level int, msg interface{}, i ...interface{}) {
	if level >= ReportingLevel {
		message := fmt.Sprint(msg)
		if level != WORKING && !strings.HasSuffix(message, "\n") {
			message += "\n"
		} else if level == WORKING && !strings.HasPrefix(message, "\r") {
			message = message + "\r"
		}
		if CurrentConfig.D.GoMonolith.ReportTimeStamp {
			log.Printf(trailTag[level]+message, i...)
		} else {
			fmt.Printf(trailTag[level]+message, i...)
		}

		// Run error handler if it exists
		if CurrentConfig.ErrorHandleFunc != nil {
			stack := string(debug.Stack())
			stackList := strings.Split(stack, "\n")
			stack = strings.Join(stackList[5:], "\n")
			go CurrentConfig.ErrorHandleFunc(level, fmt.Sprintf(fmt.Sprint(msg), i...), stack)
		}

		// Log to syslog
		if CurrentConfig.D.GoMonolith.LogTrail && level >= TrailLoggingLevel && level != WORKING {
			// Send log to syslog
			Syslogf(level, message, i...)
		}
	}
}

// TrailLoggingLevel is the minimum level to be logged into syslog
var TrailLoggingLevel = INFO
