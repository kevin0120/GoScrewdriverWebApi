package openprotocol

import (
	"context"
	"sync"
	"time"
)

type ResponseQueue struct {
	Results map[interface{}]interface{}
	mtx     sync.Mutex
}

func (q *ResponseQueue) Add(key interface{}, msg interface{}) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[key] = msg
}

func (q *ResponseQueue) update(key interface{}, msg interface{}) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	_, e := q.Results[key]
	if e {
		q.Results[key] = msg
	}
}

func (q *ResponseQueue) remove(key interface{}) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	delete(q.Results, key)
}

func (q *ResponseQueue) get(key interface{}) interface{} {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	m, e := q.Results[key]
	if e {
		return m
	} else {
		return nil
	}
}

func (q *ResponseQueue) Get(key interface{}, ctx context.Context) interface{} {
	ch := make(chan interface{})

	go func() {
		for {
			in := q.get(key)
			if in != nil {
				ch <- in
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case resp := <-ch:
		return resp
	case <-ctx.Done():
		return nil
	}

}
