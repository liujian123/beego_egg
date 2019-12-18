package models

import (
	"sync"
)

type Bucket struct {
	index   int
	rWLock  sync.RWMutex
	id2Conn map[uint64]*WsConnection
}

var (
	bucketWsConns int = 1000
)

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2Conn: make(map[uint64]*WsConnection, bucketWsConns),
	}
	return
}

func (bucket *Bucket) AddConn(wsConnection *WsConnection) {
	bucket.rWLock.Lock()
	defer bucket.rWLock.Unlock()
	bucket.id2Conn[wsConnection.ConnId] = wsConnection
}

func (bucket *Bucket) PushAll(pushJob *PushJob) {
	bucket.rWLock.RLock()
	defer bucket.rWLock.RUnlock()
	for _, WsConn := range bucket.id2Conn {
		WsConn.SendMessage(pushJob)
	}
}
