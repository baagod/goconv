package snowflake

import (
	"sync"
	"time"
)

var (
	Epoch    int64 = 1445126400000 // 2015-10-18 UTC
	SeqBits        = 12            // 序列号占用的 bit 位数
	NodeBits       = 10            // 机器号占用的 bit 位数

	timeShift       = NodeBits + SeqBits // 时间戳的偏移量
	seqMax    int64 = 1<<SeqBits - 1     // 最大序列号
)

type Snowflake struct {
	mu        sync.Mutex
	node      int64
	seq       int64
	timestamp int64
}

func New(n int64) *Snowflake {
	return &Snowflake{node: n}
}

func (a *Snowflake) Next() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now().UnixMilli() - Epoch
	if a.timestamp > now {
		now = a.timestamp + (a.timestamp - now)
	}

	if now == a.timestamp {
		a.seq = (a.seq + 1) & seqMax
		if a.seq == 0 {
			for now <= a.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		a.seq = 0
	}

	a.timestamp = now
	seq := (a.timestamp << timeShift) | (a.node << NodeBits) | a.seq

	return seq
}

func (a *Snowflake) Time(seq int64) time.Time {
	return time.UnixMilli(Epoch + (seq >> 22))
}
