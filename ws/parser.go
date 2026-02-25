package ws

type parser struct {
	buffer []byte
	status parseStatus
}

func newParser() *parser {
	return &parser{
		buffer: make([]byte, 4096),
		status: parseStatusInfo,
	}
}

type parseStatus int

const (
	parseStatusInfo parseStatus = iota + 1
	parseStatusLength
	parseStatusMask
	parseStatusPayload
)

func (p *parser) parse(data []byte) {
	// status continuing from a partial message

	// new message
}
