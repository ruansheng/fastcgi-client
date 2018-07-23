package gofastcgi

// 8 byte
type header struct {
	Version       uint8  // fastcgi protocol version
	Type          uint8  // op type
	RequestId     uint16 // request id
	ContentLength uint16 // content length
	PaddingLength uint8  // padding byte length
	Reserved      uint8  // reserve byte
}

func (h *header) init(t uint8, requestId uint16, contentLength uint16) {
	h.Version = 1
	h.Type = t
	h.RequestId = requestId
	h.ContentLength = contentLength
	h.PaddingLength = uint8(-contentLength & 7) // example: -2 & 7 = 6  supplement content enough 8 byte multiple
}
