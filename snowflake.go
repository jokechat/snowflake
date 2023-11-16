package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10
	workerMax  int64 = -1 ^ (-1 << workerBits)

	stepBits uint8 = 12
	stepMax  int64 = -1 ^ (-1 << stepBits)

	workerShift uint8 = stepBits
	timeShift   uint8 = workerBits + stepBits
)

type Worker struct {
	mu sync.Mutex

	timestamp int64     // 记录时间戳
	workerId  int64     // 当前工作节点ID
	step      int64     // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	epoch     time.Time // 开始运行时间

}

func NewWorkerWithOpts(opts ...Option) (*Worker, error) {
	w := &Worker{}
	for _, opt := range opts {
		opt.apply(w)
	}

	if w.epoch.IsZero() {
		return nil, errors.New("epoch is required")
	}

	if w.workerId < 0 || w.workerId > workerMax {
		return nil, errors.New(fmt.Sprintf("worker ID must be in [0,%d]", workerMax))
	}

	return w, nil
}

func (w *Worker) Next() ID {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixMilli()

	// 处理时钟回拨问题或处于同一毫秒 ntp 时间变化
	if now < w.timestamp {
		now = w.timestamp + 1
	}

	if now == w.timestamp {
		w.step = (w.step + 1) & stepMax
		if w.step == 0 { // 当前毫秒生成的ID已超上限,等待下一毫秒
			for now <= w.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		w.step = 0
		// todo 持久化存储当前时间,文件或者redis
	}
	w.timestamp = now
	id := (now-w.epoch.UnixMilli())<<timeShift | (w.workerId << workerShift) | w.step
	return ID(id)
}
