package writer

import "Webphish/internal"

type ApplicationOutputWriter interface {
	Capture(msg string, ready chan bool) internal.ErrorCode
}
