package metrics

import (
	"context"
	"fmt"

	"goa.design/clue/log"
)

// clue/log to prometheus logger adapter.
type logger struct {
	context.Context
}

// Implements the promhttp.Logger interface.
func (l logger) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	log.Printf(l, msg)
}
