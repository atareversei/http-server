package ringbuf

import "errors"

var ErrShortRead = errors.New("short read")
var ErrEmptyBuffer = errors.New("empty buffer")
var ErrFullBuffer = errors.New("full buffer")

type Buffer struct {
	readPos  int
	writePos int
	cap      int
	len      int
	buffer   []byte
}

func New(size int) *Buffer {
	return &Buffer{
		readPos:  0,
		writePos: 0,
		cap:      size,
		len:      0,
		buffer:   make([]byte, size),
	}
}

func (b *Buffer) IsEmpty() bool {
	return b.len == 0
}

func (b *Buffer) IsFull() bool {
	return b.cap == b.len
}

func (b *Buffer) Read(data []byte) (int, error) {
	if b.IsEmpty() {
		return 0, ErrEmptyBuffer
	}

	toRead := min(len(data), b.len)
	if b.readPos < b.writePos {
		copy(data, b.buffer[b.readPos:b.readPos+toRead])
	} else {
		firstPart := b.cap - b.readPos
		if firstPart >= toRead {
			copy(data, b.buffer[b.readPos:b.readPos+toRead])
		} else {
			copy(data, b.buffer[b.readPos:])
			copy(data[firstPart:], b.buffer[:toRead-firstPart])
		}
	}

	b.readPos = (b.readPos + toRead) % b.cap
	b.len -= toRead

	return toRead, nil
}

func (b *Buffer) ReadN(n int) ([]byte, error) {
	if n < b.len {
		return nil, ErrShortRead
	}

	data := make([]byte, n)
	_, err := b.Read(data)
	return data, err
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.IsEmpty() {
		return 0, ErrEmptyBuffer
	}

	val := b.buffer[b.readPos]
	b.readPos = (b.readPos + 1) % b.cap
	b.len--

	return val, nil
}

func (b *Buffer) GetReadBuffer() []byte {
	if b.IsEmpty() {
		return nil
	}

	if b.readPos < b.writePos {
		return b.buffer[b.readPos:b.writePos]
	} else {
		// TODO: check `b.buffer[b.readPos:b.len]`
		// It should not return the uninitialized cells of slice.
		return b.buffer[b.readPos:b.len]
	}
}

func (b *Buffer) CommitRead(n int) error {
	if n > b.len {
		return ErrShortRead
	}

	b.readPos = (b.readPos + 1) % b.cap
	b.len -= n

	return nil
}

func (b *Buffer) Write(data []byte) (int, error) {
	toWrite := len(data)
	if toWrite > (b.cap - b.len) {
		return 0, ErrFullBuffer
	}

	if b.writePos+toWrite <= b.cap {
		copy(b.buffer[b.writePos:], data)
	} else {
		firstPart := b.cap - b.writePos
		if firstPart <= toWrite {
			copy(b.buffer[b.writePos:], data)
		} else {
			copy(b.buffer[b.writePos:], data[:firstPart])
			copy(b.buffer, data[firstPart:])
		}
	}

	b.writePos = (b.writePos + toWrite) % b.cap
	b.len += toWrite

	return toWrite, nil
}

func (b *Buffer) WriteByte(data byte) error {
	if b.IsFull() {
		return ErrFullBuffer
	}

	b.buffer[b.writePos] = data
	b.writePos = (b.writePos + 1) % b.cap
	b.len++

	return nil
}

func (b *Buffer) GetWriteBuffer() []byte {
	if b.IsFull() {
		return nil
	}

	if b.writePos < b.readPos {
		return b.buffer[b.writePos:b.readPos]
	} else {
		return b.buffer[b.writePos:]
	}
}

func (b *Buffer) CommitWrite(n int) error {
	if n > b.cap-b.len {
		return ErrFullBuffer
	}

	b.writePos = (b.writePos + 1) % b.cap
	b.len += n

	return nil
}
