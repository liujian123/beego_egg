package models

import (
	"github.com/astaxie/beego/logs"
)

type PushJob struct {
	MsgType int
	MsgData []byte
}
type ConnMgr struct {
	Buckets      []*Bucket
	JobChan      []chan *PushJob
	DispatchChan chan *PushJob
}

var (
	G_ConnMgr  *ConnMgr
	bucketNums int = 2
)

func init() {
	InitConnMgr()
}

func InitConnMgr() {
	var (
		connMgr *ConnMgr
	)
	connMgr = &ConnMgr{
		Buckets:      make([]*Bucket, bucketNums),
		JobChan:      make([]chan *PushJob, bucketNums),
		DispatchChan: make(chan *PushJob, bucketNums),
	}
	for bucketIdx := 0; bucketIdx < bucketNums; bucketIdx++ {
		connMgr.Buckets[bucketIdx] = InitBucket(bucketIdx)
		connMgr.JobChan[bucketIdx] = make(chan *PushJob, 1000)
		//用于将消息分发到桶里的各个链接
		for idx := 0; idx < 32; idx++ {
			go connMgr.jobWorkerMain(bucketIdx)
		}
	}

	//用于将消息扇到多个桶中
	for idx := 0; idx < 32; idx++ {
		go connMgr.dispatchWorkerMain()
	}

	G_ConnMgr = connMgr
}

func (connMgr *ConnMgr) jobWorkerMain(bucketIdx int) {
	var (
		bucket     *Bucket = connMgr.Buckets[bucketIdx]
		pushJob *PushJob
	)

	for {
		select {
		case pushJob = <-connMgr.JobChan[bucketIdx]:
			bucket.PushAll(pushJob)
		}
	}
}

func (connMgr *ConnMgr) dispatchWorkerMain() {
	var (
		pushJob *PushJob
	)
	for {
		select {
		case pushJob = <-connMgr.DispatchChan:
			for bucketIdx, _ := range connMgr.Buckets {
				connMgr.JobChan[bucketIdx] <- pushJob
			}
		}
	}
}

func (connMgr *ConnMgr) PushAll(msg []byte) {
	var (
		pushJob *PushJob
	)
	pushJob = &PushJob{
		MsgType: 1,
		MsgData: msg,
	}

	select {
	case connMgr.DispatchChan <- pushJob:
	default:
		logs.Info("DispatchChan channel 已满")
	}
}

func (connMgr *ConnMgr) GetBucket(wsConnection *WsConnection) (bucket *Bucket) {
	bucket = connMgr.Buckets[wsConnection.ConnId%uint64(len(connMgr.Buckets))]
	return
}

func (connMgr *ConnMgr) AddConn(wsConnection *WsConnection) {
	var (
		bucket *Bucket
	)
	bucket = connMgr.GetBucket(wsConnection)
	bucket.AddConn(wsConnection)
}
