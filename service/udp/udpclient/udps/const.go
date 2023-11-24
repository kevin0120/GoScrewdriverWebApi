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
	rid          int
	buffCacheMap map[int][]byte
	futureData   *FutureData
	futureChan   chan *FutureData
	Head         string
}

func NewFuturePack(rid int) *SyncUdpFuturePack {
	return &SyncUdpFuturePack{
		rid:        rid,
		futureData: &FutureData{},
		futureChan: make(chan *FutureData, 1),
	}
}

func (c *SyncUdpFuturePack) add() {
	return
}

func (c *SyncUdpFuturePack) toArray() {
	return
}
func (c *SyncUdpFuturePack) setResult() {
	c.futureChan <- c.futureData
	close(c.futureChan)
	return
}
func (c *SyncUdpFuturePack) Result(t time.Duration) *FutureData {
	timer := time.NewTimer(t)
	for {
		select {
		case f := <-c.futureChan:
			return f
		case <-timer.C:
			c.TimeOut()
		}
	}
}
func (c *SyncUdpFuturePack) close() {
	return
}
func (c *SyncUdpFuturePack) TimeOut() {
	c.futureChan <- c.futureData
	close(c.futureChan)
	return
}
func (c *SyncUdpFuturePack) setExcept() {
	c.futureChan <- c.futureData
	close(c.futureChan)
	return
}
