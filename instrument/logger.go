package instrument

import (
	"context"
	"fmt"

	"github.com/goadesign/clue/log"
)

// micro/log to prometheus logger adapter.
type logger struct {
	context.Context
}

// Implements the promhttp.Logger interface.
func (l logger) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	log.Error(l, msg)
}
