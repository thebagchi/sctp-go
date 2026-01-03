package bitbuffer

import (
	"errors"
	"slices"
)

// InitialBufferSize is the initial capacity for the buffer in CreateWriter.
var InitialBufferSize = 64

// Codec manages a bit stream for encoding and decoding.
type Codec struct {
	Buff        []byte
	offset      uint8  // Current bit offset in the current byte (0–7). 0 = byte-aligned.
	bitsWritten uint64 // Total bits written
	bitsRead    uint64 // Total bits read
}

// CreateWriter creates a new Codec for writing.
func CreateWriter() *Codec {
	return &Codec{
		Buff: make([]byte, 0, InitialBufferSize),
	}
}

// CreateReader creates a new Codec for reading from existing data.
func CreateReader(data []byte) *Codec {
	return &Codec{
		Buff:   data,
		offset: 0,
	}
}

// Len returns the number of complete bytes currently in the buffer.
func (c *Codec) Len() int {
	return len(c.Buff)
}

// Cap returns the capacity of the underlying buffer.
func (c *Codec) Cap() int {
	return cap(c.Buff)
}

// NumWritten returns the total number of bits written.
func (c *Codec) NumWritten() uint64 {
	return c.bitsWritten
}

// NumRead returns the total number of bits read.
func (c *Codec) NumRead() uint64 {
	return c.bitsRead
}

// Bytes returns the encoded data trimmed to the exact number of bytes needed.
// Includes the partial final byte if bitsWritten is not a multiple of 8.
func (c *Codec) Bytes() []byte {
	if c.bitsWritten == 0 {
		return nil
	}
	return c.Buff
}

// grow ensures space for at least n more bytes.
func (c *Codec) grow(n int) {
	if cap(c.Buff) < len(c.Buff)+n {
		capacity := max(cap(c.Buff)*2, len(c.Buff)+n)
		c.Buff = slices.Grow(c.Buff, capacity-len(c.Buff))
	}
	c.Buff = c.Buff[:len(c.Buff)+n]
}

// append appends a new zero byte and resets offset to 0.
func (c *Codec) append() {
	c.grow(1)
	c.offset = 0
}

// incrementRead increments the bits read counter.
func (c *Codec) incrementRead(bits uint64) {
	c.bitsRead += bits
}

// incrementWrite increments the bits written counter.
func (c *Codec) incrementWrite(bits uint64) {
	c.bitsWritten += bits
}

// Write writes the least significant 'num' bits of value (1 ≤ num ≤ 64).
func (c *Codec) Write(num uint8, value uint64) error {
	if num == 0 || num > 64 {
		return errors.New("bit count must be between 1 and 64")
	}

	value &= (1 << num) - 1

	bitsLeft := num
	for bitsLeft > 0 {
		if c.offset == 8 {
			c.append()
		}
		if len(c.Buff) == 0 {
			c.append()
		}

		var (
			remaining = uint8(8 - c.offset)
			writeNow  = min(bitsLeft, remaining)
			shift     = bitsLeft - writeNow
			bits      = uint8(value>>shift) & ((1 << writeNow) - 1)
			bitShift  = remaining - writeNow
			idx       = len(c.Buff) - 1
		)

		c.Buff[idx] = c.Buff[idx] | (bits << bitShift)

		c.offset += writeNow
		bitsLeft -= writeNow
	}

	c.incrementWrite(uint64(num))
	return nil
}

// Read reads the next num bits from the bit stream, returning them as a uint64.
func (c *Codec) Read(num uint8) (uint64, error) {
	if num == 0 {
		return 0, nil
	}
	if num > 64 {
		return 0, errors.New("bit count must be between 1 and 64")
	}

	if c.Len() == 0 {
		return 0, errors.New("no more data")
	}

	var result uint64
	bitsLeft := num

	for bitsLeft > 0 {
		if c.offset == 8 {
			c.Buff = c.Buff[1:]
			c.offset = 0
			if len(c.Buff) == 0 {
				return 0, errors.New("unexpected end of data")
			}
		}

		var (
			remaining = uint8(8 - c.offset)
			readNow   = min(bitsLeft, remaining)
			mask      = uint8((1 << readNow) - 1)
			shift     = remaining - readNow
			bits      = uint64((c.Buff[0] >> shift) & mask)
		)

		result = (result << readNow) | bits

		c.offset += readNow
		bitsLeft -= readNow
	}

	c.incrementRead(uint64(num))
	return result, nil
}

// WriteBytes writes full octets continuing from the current bit offset.
// Equivalent to repeated Write(8, uint64(b)) for each byte.
// Does NOT force alignment — caller must Align() if required (e.g., APER octet string contents).
func (c *Codec) WriteBytes(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Fast path: already byte-aligned
	if c.offset == 0 {
		length := len(c.Buff)
		c.grow(len(data))
		copy(c.Buff[length:], data)
		c.incrementWrite(uint64(len(data) * 8))
		return nil
	}

	// Slow path: pack each byte using general Write
	for _, b := range data {
		if err := c.Write(8, uint64(b)); err != nil {
			return err
		}
	}
	return nil
}

// ReadBytes reads exactly n full octets (bytes) from the bit stream.
// Continues from current bit offset. Equivalent to calling Read(8) n times.
// Returns error if insufficient data.
func (c *Codec) ReadBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("negative byte count")
	}
	if n == 0 {
		return []byte{}, nil
	}

	// Fast path: already byte-aligned
	if c.offset == 0 {
		if len(c.Buff) < n {
			return nil, errors.New("insufficient data")
		}
		result := make([]byte, n)
		copy(result, c.Buff[:n])
		c.Buff = c.Buff[n:]
		c.incrementRead(uint64(n * 8))
		return result, nil
	}

	// Slow path: read each byte using general Read
	result := make([]byte, n)
	for i := range result {
		val, err := c.Read(8)
		if err != nil {
			return nil, err
		}
		result[i] = uint8(val)
	}
	return result, nil
}

// Align advances to the next byte boundary by appending a new byte if necessary.
// Unused bits in the previous byte remain zero.
// Used explicitly by the PER encoder when alignment is required (e.g., APER).
func (c *Codec) Align() error {
	if c.offset > 0 && c.offset < 8 {
		// Partial byte - need to pad and move to next byte
		c.append()
	} else if c.offset == 8 {
		// Just finished a byte - just reset offset
		c.offset = 0
	}
	// If offset == 0, already aligned, do nothing
	return nil
}

// Advance skips remaining bits to reach the next byte boundary (for reading).
// Does not modify the buffer, just updates offset and advances Buff if offset reaches 8.
func (c *Codec) Advance() error {
	if c.offset > 0 {
		c.offset = 8
		// Trigger the automatic Buff advancement
		if len(c.Buff) > 0 {
			c.Buff = c.Buff[1:]
			c.offset = 0
		}
	}
	return nil
}
