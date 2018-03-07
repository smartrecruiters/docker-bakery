package commons

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	// Debug is true when the DEBUG env var is not empty
	IsDebugEnabled = os.Getenv("DEBUG") != ""
)

type MessageProvider func() string

func Debug(message string) {
	printfMsg(message)
}

func Debugf(message string, args ...interface{}) {
	printfMsg(message, args...)
}

// Lazy debug should be used when message construction is heavy and should be executed only when debug is enabled
func LazyDebug(getMsgFn MessageProvider) {
	if IsDebugEnabled {
		printfMsg(getMsgFn())
	}
}

func printfMsg(msg string, v ...interface{}) {
	if IsDebugEnabled {
		log.Print(color.CyanString(fmt.Sprintf(msg, v...)))
	}
}
