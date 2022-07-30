package utils

import (
	"strconv"
	"sync"
	"time"
)

//	定义基本变量，雪花算法,时间戳，机器id，以及二者偏移量
const (
	twepoch      = int64(1483228800000) //开始时间戳(根据实际情况考虑)
	workidBits   = uint(4)              // 机器id所占用位数
	sequenceBits = uint(18)             // 序列所占位数

	workeridMax    = int64(-1 ^ (-1 << workidBits))   //最大机器id
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //最大序列数
	workIdShift    = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workidBits
)

// SnowFlake 雪花算法所需要的结构
type SnowFlake struct {
	sync.Mutex
	timestamp int64
	workerId  int64
	sequence  int64
}

func NewSnowflake() *SnowFlake {
	return &SnowFlake{
		timestamp: 0,
		workerId:  int64(GenerateRandomNumber(0, 1023)),
		sequence:  0,
	}
}

func (s *SnowFlake) Generate() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}
	s.timestamp = now
	r := (now-0)<<timestampShift | (s.workerId << workIdShift) | (s.sequence)
	s.Unlock()
	return r
}

func GetUniqLogID() string {
	snow := NewSnowflake()
	rString := RandomString(5)
	return rString + strconv.FormatInt(snow.Generate(), 10)
}
