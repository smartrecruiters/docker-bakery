package commons

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	// Debug is true when the DEBUG env var is not empty.
	IsDebugEnabled = os.Getenv("DEBUG") != ""
)

// Allows for definition of MessageProvider that will be invoked in order to obtain the message before logging it.
type MessageProvider func() string

// Prints message when debug mode is enabled.
func Debug(message string) {
	printfMsg(message)
}

// Prints message when debug mode is enabled. Substitutes format with provided arguments. Works like fmt.Sprintf.
func Debugf(message string, args ...interface{}) {
	printfMsg(message, args...)
}

// Lazy debug should be used when message construction is heavy and should be executed only when debug is enabled
func LazyDebug(getMsgFn MessageProvider) {
	if IsDebugEnabled {
		printfMsg(getMsgFn())
	}
}

// Prints the message if logging is enabled.
func printfMsg(msg string, v ...interface{}) {
	if IsDebugEnabled {
		log.Print(color.CyanString(fmt.Sprintf(msg, v...)))
	}
}
