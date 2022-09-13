//go:build debug

package logger

import (
	"log"
)

func Fatalf(s string, args ...interface{}) {
	log.Fatalf(s, args)
}
