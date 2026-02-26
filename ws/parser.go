package ws

import "github.com/atareversei/http-server/internal/ringbuf"

const bufferSize = 65536

type parser struct {
	buffer       ringbuf.Buffer
	status       parseStatus
	currentFrame Frame
}

func newParser() *parser {
	return &parser{
		buffer: *ringbuf.New(bufferSize),
		status: parseInitialized,
	}
}

func (p *parser) parse() (Frame, error) {
	if p.status == parseInitialized {
		p.parseInfo()
	}
	if p.status == infoParsed {
		p.parseLength()
	}
	if p.status == lengthParsed {
		p.parseMask()
	}
	if p.status == maskParsed {
		p.parsePayload()
	}
	if p.status == payloadParsed {
		p.cleanUp()
	}

	return p.currentFrame, nil
}

func (p *parser) parseInfo() {
	if p.status > parseInitialized {
		return
	}

	headers, err := p.buffer.ReadN(2)
	if err != nil {
	}

	b1 := headers[0]
	b2 := headers[1]

	finl := b1 & (0b10000000) >> 7
	rsv1 := b1 & (0b01000000) >> 6
	rsv2 := b1 & (0b00100000) >> 5
	rsv3 := b1 & (0b00010000) >> 4
	opcd := b1 & (0b00001111)

	mask := b2 & (0b10000000) >> 7
	leng := b2 & (0b01111111)

	p.currentFrame.isFin = finl == 1
	p.currentFrame.rsv1 = rsv1
	p.currentFrame.rsv2 = rsv2
	p.currentFrame.rsv3 = rsv3
	p.currentFrame.opCode = byteToOpCode(opcd)

	p.currentFrame.hasMask = mask == 1
	if p.currentFrame.hasMask {
		p.currentFrame.maskByteLeft = 4
	}

	p.currentFrame.payloadType = byteToPayloadType(leng)
	if p.currentFrame.payloadType == shortPayload {
		p.currentFrame.payloadByteLeft = int(leng)
	}

	p.status = infoParsed
}

func (p *parser) parseLength() {
	p.status = lengthParsed
}

func (p *parser) parseMask() {
	p.status = maskParsed
}

func (p *parser) parsePayload() {
	p.status = payloadParsed
}

func (p *parser) cleanUp() {
	p.status = parseInitialized
}

type parseStatus int

const (
	parseInitialized parseStatus = iota + 1
	infoParsed
	lengthParsed
	maskParsed
	payloadParsed
)
