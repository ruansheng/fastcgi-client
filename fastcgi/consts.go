package fastcgi

const (
	FcgiBeginRequest uint8 = iota + 1
	FCGI_ABORT_REQUEST  // 2
	FCGI_END_REQUEST    // 3
	FcgiParams          // 4
	FcgiStdin			// 5
	FCGI_STDOUT
	FCGI_STDERR
	FCGI_DATA
	FCGI_GET_VALUES
	FCGI_GET_VALUES_RESULT
	FCGI_UNKNOWN_TYPE
	FCGI_MAXTYPE = FCGI_UNKNOWN_TYPE
)

const (
	FCGI_RESPONDER uint8 = iota + 1
	FCGI_AUTHORIZER
	FCGI_FILTER
)

const (
	maxWrite = 6553500 // maximum record body
	maxPad   = 255
)
