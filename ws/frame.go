package ws

import (
	"encoding/binary"
	"errors"
)

// TODO: enable streaming for continuation frames with large payloads
type frame struct {
	isFin          bool
	rsv1           byte
	rsv2           byte
	rsv3           byte
	opCode         opCode
	payloadLenType payloadLenType
	hasMask        bool
	mask           []byte
	len            uint64
	content        []byte
}

func newFrame() frame {
	return frame{}
}

type opCode byte

const (
	opContinuation opCode = iota
	opText
	opBinary

	opRsv1
	opRsv2
	opRsv3
	opRsv4
	opRsv5

	opClose
	opPing
	opPong

	opRsvA
	opRsvB
	opRsvC
	opRsvD
	opRsvE

	opInvalid
)

func byteToOpCode(opCode byte) opCode {
	switch opCode {
	case 0b0000:
		return opContinuation
	case 0b0001:
		return opText
	case 0b0010:
		return opBinary

	case 0b0011:
		return opRsv1
	case 0b0100:
		return opRsv2
	case 0b0101:
		return opRsv3
	case 0b0110:
		return opRsv4
	case 0b0111:
		return opRsv5

	case 0b1000:
		return opClose
	case 0b1001:
		return opPing
	case 0b1010:
		return opPong

	case 0b1011:
		return opRsvA
	case 0b1100:
		return opRsvB
	case 0b1101:
		return opRsvC
	case 0b1110:
		return opRsvD
	case 0b1111:
		return opRsvE

	default:
		return opInvalid
	}
}

type payloadLenType byte

const (
	shortPayloadLen payloadLenType = iota
	mediumPayloadLen
	extendedPayloadLen
	invalidPayloadLen
)

func byteToPayloadLenType(payload byte) payloadLenType {
	if payload < 126 {
		return shortPayloadLen
	}

	if payload == 126 {
		return mediumPayloadLen
	}

	if payload == 127 {
		return extendedPayloadLen
	}

	return invalidPayloadLen
}

func bytesToLen(bytes []byte) (uint64, error) {
	if len(bytes) > 8 {
		return 0, errors.New("bytes slice is invalid")
	}
	padded := make([]byte, 8)
	copy(padded[8-len(bytes):], bytes)

	payload := binary.BigEndian.Uint64(padded)

	return payload, nil
}
