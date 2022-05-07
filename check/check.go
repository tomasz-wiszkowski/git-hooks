package check

import (
	"fmt"
	"log"
)

// Given an error and a message text, emit panic() if the error is not nil
// reporting the (formatted) error message.
// If error is nil no action is taken.
func Err(err error, text string, args ...interface{}) {
	if err != nil {
		s := fmt.Sprintf(text, args...)
		log.Output(2, "Assertion failure")
		log.Panicf("%s: %s", s, err)
	}
}

// Given condition and a message text, emit panic() if the condition
// evaluates to false, reporting the (formatted) error message.
// If the condition evaluates to true no action is taken.
func True(condition bool, text string, args ...interface{}) {
	if !condition {
		log.Output(2, "Assertion failure")
		log.Panicf(text, args...)
	}
}
