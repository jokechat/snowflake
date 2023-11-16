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
