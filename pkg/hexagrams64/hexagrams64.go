package hexagrams64

import (
	"fmt"
)

/*
 * Encodings
 */
type Encoding struct {
	padChar rune
}

// 没有NewEncoding，因为encoder alphabet是固定的。

// padding只能是3-byte rune，即[1<<11, 1<<16)，并且不在_64范围内。
// 否则必须是_NoPadding，其它情况均为无效值。
// 也许还应该限制必须是可打印字符	unicode.IsPrint(padding)
func (enc Encoding) WithPadding(padding rune) Encoding {
	if padding == enc.padChar {
		return enc
	}

	switch {
	case padding == NoPadding:
	case 1<<11 <= padding && padding < '\u4DC0': // < _64
	case '\u4DFF' < padding && padding < 1<<16: // > _64
	default:
		panic("invalid padding")
	}
	return Encoding{padChar: padding}
}

// 没有Strict，因为base64是以byte为单位，可以选择是否忽略\r \n等byte。而这里是以3-byte rune为单位，无法处理单字节值

/*
 * Encoder
 */

func encode3(b0, b1, b2 byte) (u3, u2, u1, u0 uint8) {
	val := uint32(b0)<<16 | uint32(b1)<<8 | uint32(b2)
	u3 = uint8(val >> 18 & 0x3F)
	u2 = uint8(val >> 12 & 0x3F)
	u1 = uint8(val >> 6 & 0x3F)
	u0 = uint8(val & 0x3F)
	return
}

func decode4(u3, u2, u1, u0 uint8) (b0, b1, b2 byte) {
	val := uint32(u3)<<26 | uint32(u2)<<20 | uint32(u1)<<14 | uint32(u0)<<8
	b0 = byte(val >> 24)
	b1 = byte(val >> 16)
	b2 = byte(val >> 8)
	return
}

func (enc Encoding) Encode(dst, src []byte) {
	di, si := 0, 0
	for n := (len(src) / 3) * 3; si < n; si += 3 {
		// Convert 3x 8bit source bytes into 4 runes
		u3, u2, u1, u0 := encode3(src[si+0], src[si+1], src[si+2])
		r3, r2, r1, r0 := table[u3], table[u2], table[u1], table[u0]
		di = append4rune(dst, di, r3, r2, r1, r0)
	}

	switch remain := len(src) - si; remain {
	case 0:
		return
	case 1: // 2 padding
		u3, u2, _, _ := encode3(src[si+0], 0, 0) // 不足位补0
		r3, r2, _, _ := table[u3], table[u2], enc.padChar, enc.padChar
		di = appendrune(dst, di, r3)
		di = appendrune(dst, di, r2)
		if enc.padChar != NoPadding {
			di = appendrune(dst, di, enc.padChar)
			di = appendrune(dst, di, enc.padChar)
		}
	case 2: // 1 padding
		u3, u2, u1, _ := encode3(src[si+0], src[si+1], 0) // 不足位补0
		r3, r2, r1, _ := table[u3], table[u2], table[u1], enc.padChar
		di = appendrune(dst, di, r3)
		di = appendrune(dst, di, r2)
		di = appendrune(dst, di, r1)
		if enc.padChar != NoPadding {
			di = appendrune(dst, di, enc.padChar)
		}
	}
}

func append4rune(dst []byte, di int, r3, r2, r1, r0 rune) int {
	dst12 := dst[di : di+12]
	dst12[0] = 0xE0 | byte(r3>>12)
	dst12[1] = 0x80 | byte(r3>>6)&0x3F
	dst12[2] = 0x80 | byte(r3)&0x3F
	dst12[3] = 0xE0 | byte(r2>>12)
	dst12[4] = 0x80 | byte(r2>>6)&0x3F
	dst12[5] = 0x80 | byte(r2)&0x3F
	dst12[6] = 0xE0 | byte(r1>>12)
	dst12[7] = 0x80 | byte(r1>>6)&0x3F
	dst12[8] = 0x80 | byte(r1)&0x3F
	dst12[9] = 0xE0 | byte(r0>>12)
	dst12[10] = 0x80 | byte(r0>>6)&0x3F
	dst12[11] = 0x80 | byte(r0)&0x3F
	return di + 12
}

func appendrune(dst []byte, di int, r0 rune) int {
	dst3 := dst[di : di+3]
	dst3[0] = 0xE0 | byte(r0>>12)
	dst3[1] = 0x80 | byte(r0>>6)&0x3F
	dst3[2] = 0x80 | byte(r0)&0x3F
	return di + 3
}

func (enc *Encoding) AppendEncode(dst, src []byte) []byte {
	n := enc.EncodedLen(len(src))
	if n -= cap(dst) - len(dst); n > 0 {
		dst = append(dst[:cap(dst)], make([]byte, n)...)[:len(dst)]
	}
	enc.Encode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n]
}

func (enc Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}

func (enc Encoding) EncodedLen(n int) int {
	if enc.padChar == NoPadding {
		// NoPadding:	3n->12n		3n+1->4n+2->12n+6		3n+2->4n+3->12n+9
		return n/3*12 + (n%3*8+5)/6*3
	}
	// Padding:		3n->12n		3n+1->3(n+1)->12(n+1)		3n+2->3(n+1)->12(n+1)
	return (n + 2) / 3 * 12
}

/*
 * Decoder
 */

func (enc Encoding) decodeQuantum(dst, src []byte, si int) (nsi, n int, err error) {
	// Decode quantum using the base64 alphabet
	var dbuf [4]byte
	dlen := 4

	for j := 0; j < len(dbuf); j++ {
		if len(src) == si {
			switch {
			case j == 0:
				return si, 0, nil
			case j == 1, enc.padChar != NoPadding:
				return si, 0, fmt.Errorf("expected NoPadding at %d", si-j*3)
			}
			dlen = j
			break
		}

		in := decoderune(src[si:])
		si += 3

		out := decodeMap(in)
		if out != 0xff { // nomal
			dbuf[j] = out
			continue
		}

		if in != enc.padChar { // out == 0xff
			return si, 0, fmt.Errorf("expected %c at %d,actual: %c", enc.padChar, si-3, in)
		}

		// in == enc.padChar
		switch j {
		case 0, 1:
			return si, 0, fmt.Errorf("incorrect padding")
		case 2:
			if si == len(src) {
				return si, 0, fmt.Errorf("not enough padding")
			}

			in = decoderune(src[si:])
			if in != enc.padChar {
				return si, 0, fmt.Errorf("expected %c at %d,actual: %c", enc.padChar, si, in)
			}
			si += 3
		}

		if si < len(src) {
			err = fmt.Errorf("trailing garbage")
		}
		dlen = j
		break
	}

	// Convert 4x 6bit source bytes into 3 bytes
	val := uint32(dbuf[0])<<18 | uint32(dbuf[1])<<12 | uint32(dbuf[2])<<6 | uint32(dbuf[3])
	dbuf[2], dbuf[1], dbuf[0] = byte(val>>0), byte(val>>8), byte(val>>16)
	switch dlen {
	case 4:
		dst[2] = dbuf[2]
		dbuf[2] = 0
		fallthrough
	case 3:
		dst[1] = dbuf[1]
		dbuf[1] = 0
		fallthrough
	case 2:
		dst[0] = dbuf[0]
	}

	return si, dlen - 1, err
}

func (enc Encoding) AppendDecode(dst, src []byte) ([]byte, error) {
	// Compute the output size without padding to avoid over allocating.
	n := len(src)
	for n > 0 && decoderune(src[n-3:]) == enc.padChar {
		n -= 3
	}
	n /= 4

	if n -= cap(dst) - len(dst); n > 0 {
		dst = append(dst[:cap(dst)], make([]byte, n)...)[:len(dst)]
	}

	n, err := enc.Decode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n], err
}

func (enc Encoding) DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

func (enc Encoding) Decode(dst, src []byte) (n int, err error) {
	if len(src)%3 != 0 {
		return 0, fmt.Errorf("invalid len %d", len(src))
	}

	si := 0
	for len(src)-si >= 12 && len(dst)-n >= 4 {
		r3, r2, r1, r0 := decode4rune(src[si:])
		u3 := decodeMap(r3)
		u2 := decodeMap(r2)
		u1 := decodeMap(r1)
		u0 := decodeMap(r0)
		if u3|u2|u1|u0 == 0xFF {
			var ninc int
			si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
			n += ninc
			if err != nil {
				return n, err
			}
		} else {
			b0, b1, b2 := decode4(u3, u2, u1, u0)
			dst4 := dst[n : n+4]
			dst4[0] = b0
			dst4[1] = b1
			dst4[2] = b2
			n += 3
			si += 12
		}
	}

	for si < len(src) {
		var ninc int
		si, ninc, err = enc.decodeQuantum(dst[n:], src, si)
		n += ninc
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func (enc Encoding) DecodedLen(n int) int {
	// Padding:		3n->12n->3n		3n+1->12n+12->3n+3		3n+2->12n+12->3n+3
	// NoPadding:	3n->12n->3n		3n+1->12n+6->3n+1		3n+2->12n+9->3n+2
	return n / 4
}

func decoderune(src []byte) rune {
	src3 := src[0:3]
	return rune(src3[0]&0x0F)<<12 | rune(src3[1]&0x3F)<<6 | rune(src3[2]&0x3F)
}

func decode4rune(src []byte) (r3, r2, r1, r0 rune) {
	src12 := src[0:12]
	r3 = rune(src12[0]&0x0F)<<12 | rune(src12[1]&0x3F)<<6 | rune(src12[2]&0x3F)
	r2 = rune(src12[3]&0x0F)<<12 | rune(src12[4]&0x3F)<<6 | rune(src12[5]&0x3F)
	r1 = rune(src12[6]&0x0F)<<12 | rune(src12[7]&0x3F)<<6 | rune(src12[8]&0x3F)
	r0 = rune(src12[9]&0x0F)<<12 | rune(src12[10]&0x3F)<<6 | rune(src12[11]&0x3F)

	return
}

// 由于是rune -> uint8的关系，就不能再使用[256]uint8的技巧了。
// 另外因为encoder alphabet是固定的，因此decodeMap没有必要作为Encoding的一部分
func decodeMap(r rune) uint8 {
	if '\u4DC0' <= r && r <= '\u4DFF' {
		return uint8(r - '\u4DC0')
	}
	return 0xFF
}
