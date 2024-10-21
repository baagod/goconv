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
	mu      sync.Mutex
	node    int64
	seq     int64
	elapsed int64
}

func New(node int64) *Snowflake {
	return &Snowflake{node: node}
}

func (f *Snowflake) Next() int64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now().UnixMilli() - Epoch
	if f.elapsed > now {
		now = f.elapsed + (f.elapsed - now)
	}

	if now == f.elapsed {
		f.seq = (f.seq + 1) & seqMax
		if f.seq == 0 {
			for now <= f.elapsed {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		f.seq = 0
	}

	f.elapsed = now
	seq := (f.elapsed << timeShift) | (f.node << NodeBits) | (f.seq)

	return seq
}

func (f *Snowflake) Time(seq int64) time.Time {
	return time.UnixMilli(Epoch + (seq >> 22))
}
