//go:build !debug

package logger

func Fatalf(s string, args ...interface{}) {}
