package snowflake

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"strconv"
	"strings"
	"time"
)

type ID uint64

func (i ID) Uint64() uint64 {
	return uint64(i)
}

func (i ID) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i ID) Base32() string {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, i)
	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(bytesBuffer.Bytes())
}

func shuffleBytes(b []byte) {
	n := len(b)
	for j := 0; j < n/2; j++ {
		k := (j*7 + 3) % n
		b[j], b[k] = b[k], b[j]
	}
}

func (i ID) Base62() string {
	const alphabet = "r4fVqD0jM7pS2cLTKmXxzGY5vFNuB8n1RkOZyaQsg6WCe9lhHJdPi3wUtoAbE"
	num := i.Uint64() ^ 0x9E3779B97F4A7C15
	if num == 0 {
		return string(alphabet[0])
	}
	var encoded []byte
	base := uint64(len(alphabet))
	for num > 0 {
		rem := num % base
		encoded = append(encoded, alphabet[rem])
		num /= base
	}
	for j, k := 0, len(encoded)-1; j < k; j, k = j+1, k-1 {
		encoded[j], encoded[k] = encoded[k], encoded[j]
	}
	shuffleBytes(encoded)
	return string(encoded)
}

func (i ID) Base32Lower() string {
	return strings.ToLower(i.Base32())
}

func (i ID) UnixMilli(epoch time.Time) int64 {
	return epoch.UnixMilli() + int64(i.Uint64()>>timeShift)
}

func (i ID) WorkId() uint64 {
	d := i.Uint64() >> workerShift & uint64(workerMax)
	return d
}

func (i ID) Step() uint64 {
	d := i.Uint64() & uint64(stepMax)
	return d
}

func (i ID) Time(epoch time.Time) time.Time {
	return time.UnixMilli(i.UnixMilli(epoch))
}
