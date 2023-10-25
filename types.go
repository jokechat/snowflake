package snowflake

import "time"

type Option interface {
	apply(*Worker)
}

//type options struct {
//	opts []Option
//}

type OptionFunc func(*Worker)

func (o OptionFunc) apply(worker *Worker) {
	o(worker)
}

func WithEpoch(epoch time.Time) OptionFunc {
	return func(worker *Worker) {
		worker.epoch = epoch
	}
}

func WithWorkerId(workerId int64) OptionFunc {
	return func(worker *Worker) {
		worker.workerId = workerId
	}
}
