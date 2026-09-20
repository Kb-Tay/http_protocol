package buffer

const INIT_SIZE = 8

type Buffer struct {
	buf         []byte
	size        int
	nextInd     int
	bytesParsed int // number of bytes parsed
	bytesRead   int
}

func New() *Buffer {
	b := Buffer{
		buf:         make([]byte, 0, INIT_SIZE),
		size:        INIT_SIZE,
		nextInd:     0,
		bytesParsed: 0,
		bytesRead:   0,
	}

	return &b
}

func (b *Buffer) Read(bytes []byte) {
	n := len(bytes)

	// Resizing
	if b.nextInd+n >= b.size {
		newBuf := make([]byte, 0, b.size*2)
		newBuf = append(newBuf, b.buf[:b.nextInd]...)

		b.buf = newBuf
		b.size = b.size * 2
	}

	// Adding
	b.buf = append(b.buf, bytes...)
	b.nextInd += n
	b.bytesRead += n
}

func (b *Buffer) Parsed(n int) {
	rest := b.buf[n:]
	b.buf = b.buf[:0] // emptying the buffer
	b.buf = append(b.buf, rest...)

	b.bytesParsed += n
	b.nextInd = len(rest)
}

func (b *Buffer) GetBuffer() []byte {
	return b.buf[:b.nextInd]
}
