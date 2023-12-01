package snowflake

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
	"time"
)

func TestCode(t *testing.T) {
	now := time.Now().UnixNano() / 1e6

	spew.Dump(now)
	spew.Dump(now << 22)

	tt := fmt.Sprintf("%064b", now)
	id := fmt.Sprintf("%064b", now<<22)

	spew.Dump(1023 & 1022)
	spew.Dump(1022 & 1023)
	spew.Dump(1025 & 1023)
	spew.Dump(tt)
	spew.Dump(id)
	fmt.Printf("%s %s %s %s",
		id[0:1],
		id[1:42],
		id[42:52],
		id[52:64],
	)
}

func TestId(t *testing.T) {
	// 项目启动时间
	epoch := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.Local)

	w, _ := NewWorkerWithOpts(
		WithEpoch(epoch),
		WithWorkerId(1),
	)
	now := time.Now().UnixMilli()
	step := 1

	ids := make([]ID, 0)
	for time.Now().UnixMilli() <= now+1 {
		id := w.Next()
		ids = append(ids, id)

		spew.Dump(id.String())
		spew.Dump(id.Base32())
		spew.Dump(id.Base32Lower())
		spew.Dump(id.UnixMilli(w.epoch))
		spew.Dump(id.Time(w.GetEpoch()))
		step++
	}

	//spew.Dump(ids)
	spew.Dump(step)

}

func TestId2(t *testing.T) {
	epoch := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.Local)

	w, _ := NewWorkerWithOpts(
		WithEpoch(epoch),
		WithWorkerId(16),
	)
	id := w.Next()
	id = w.Next()

	spew.Dump(id.Uint64())
}

func TestId3(t *testing.T) {
	spew.Dump(645228270933049345 & stepMax)
}
