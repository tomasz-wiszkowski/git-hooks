package log

import (
	"fmt"
	"log"
)

/// Given an error and a message text, emit panic() if the error is not nil
/// reporting the (formatted) error message.
/// If error is nil no action is taken.
func Check(err error, text string, args ...interface{}) {
	if err != nil {
		s := fmt.Sprintf(text, args...)
		log.Panicf("%s: %s", s, err)
	}
}
