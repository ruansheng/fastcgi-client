package fastcgi_client

const (
	FcgiBeginRequest uint8 = 1
	FcgiParams uint8 = 4
	FcgiStdin uint8 = 5
)

const (
	maxWrite = 6553500 // maximum record body
	maxPad   = 255
)
