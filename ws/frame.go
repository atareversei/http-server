package ws

// TODO: enable streaming for continuation frames with large payloads
type Frame struct {
	isFin  bool
	rsv1   byte
	rsv2   byte
	rsv3   byte
	opCode opCode

	hasMask         bool
	payloadType     payloadType
	payloadByteLeft int

	headerBytesLeft int
	lengthByteLeft  int
	maskByteLeft    int

	content []byte
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

type payloadType byte

const (
	shortPayload payloadType = iota
	mediumPayload
	extendedPayload
	invalidPayload
)

func byteToPayloadType(payload byte) payloadType {
	if payload >= 0 && payload < 126 {
		return shortPayload
	}

	if payload == 126 {
		return mediumPayload
	}

	if payload == 127 {
		return extendedPayload
	}

	return invalidPayload
}
