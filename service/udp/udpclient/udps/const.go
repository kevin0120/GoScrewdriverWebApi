package udps

import (
	"time"
)

const (
	PackSize        = 1400
	SinglePackBytes = "\xff\xff\xff\xff"
	StartPackBytes  = "\xff\xff"
)

const (
	SUCCESS     = 0x00
	FAILED      = 0xFF
	EXCEPTION   = 0xFE
	TIMEOUT     = 0xFD
	UNCONNECTED = 0x90
	COMPLETED   = 0x01
)

type FutureData struct {
	result  byte
	content []byte
	ex      string
	errCode string
}
type SyncUdpFuturePack struct {
	rid          int32
	buffCacheMap map[int32][]byte
	packLen      int
	futureData   *FutureData
	futureChan   chan *FutureData
	Head         []byte
}

func NewFuturePack(rid int32) *SyncUdpFuturePack {
	return &SyncUdpFuturePack{
		rid:          rid,
		packLen:      -1,
		buffCacheMap: map[int32][]byte{},
		futureData:   &FutureData{},
		futureChan:   make(chan *FutureData, 3),
		Head:         []byte(""),
	}
}

//Close 手动关闭一下,遇到连续两次相同请求时
func (c *SyncUdpFuturePack) Close() {
	c.timeOut()
	return
}

// Result 调用段,阻塞等待返回结果
func (c *SyncUdpFuturePack) Result(t time.Duration) *FutureData {
	timer := time.NewTimer(t)
	for {
		select {
		case f := <-c.futureChan:
			return f
		case <-timer.C:
			c.timeOut()
		}
	}
}

// Add 每次从下位机收到一个包,调用
func (c *SyncUdpFuturePack) Add(pid int32, data []byte) {
	if _, ok := c.buffCacheMap[pid]; !ok {
		if pid < 0 {
			pid *= -1
			c.packLen = int(pid)
		}
		c.buffCacheMap[pid] = data
		if len(c.buffCacheMap) == c.packLen {
			c.setResult()
		}
	} else {
		c.setExcept()
	}
}

// setResult 从下位机收到最后一个包,调用
func (c *SyncUdpFuturePack) setResult() {
	buff := c.Head
	// 使用for range遍历map
	for _, value := range c.buffCacheMap {
		buff = append(buff, value...)
	}
	c.futureData.result = SUCCESS
	c.futureData.content = buff
	c.futureChan <- c.futureData
	return
}

// timeOut 长时间没收到或没收全,调用  ps 当连续两次产生同意请求时,也把原来的对象timeout
func (c *SyncUdpFuturePack) timeOut() {
	c.futureData.result = FAILED
	c.futureChan <- c.futureData
	return
}

// setExcept 出现异常,调用
func (c *SyncUdpFuturePack) setExcept() {
	c.futureData.result = FAILED
	c.futureChan <- c.futureData
	return
}
