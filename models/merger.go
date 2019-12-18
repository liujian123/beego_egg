package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type PushContext struct {
	msg *json.RawMessage
}

type BizMessage struct {
	Item []*json.RawMessage `json:"items"`
}

type PushBatch struct {
	items       []*json.RawMessage
	commitTimer *time.Timer
}

type MergeWorker struct {
	contextChan chan *PushContext
	timeOutChan chan *PushBatch
	allBatch    *PushBatch
}

type Merger struct {
	broadcastWorker *MergeWorker
}

var (
	G_merger *Merger
)

func init() {
	InitMergeWorker()
}

func InitMergeWorker() {
	var (
		mergeWorker *MergeWorker
		merger      *Merger
	)
	mergeWorker = &MergeWorker{
		contextChan: make(chan *PushContext, 1000),
		timeOutChan: make(chan *PushBatch, 10),
	}
	merger = &Merger{
		broadcastWorker: mergeWorker,
	}
	go mergeWorker.MergerWorkerMain()
	G_merger = merger
}

func (mergeWorker *MergeWorker) MergerWorkerMain() {
	var (
		context      *PushContext
		batch        *PushBatch
		timeOutBatch *PushBatch
		isCreated    bool
	)
	for {
		select {
		case context = <-mergeWorker.contextChan:
			batch = mergeWorker.allBatch
			if batch == nil {
				batch = &PushBatch{}
				mergeWorker.allBatch = batch
				isCreated = true
			}
			batch.items = append(batch.items, context.msg)
			if isCreated {
				fmt.Println("isCreated::::::", isCreated)
				batch.commitTimer = time.AfterFunc(time.Duration(time.Second*1), mergeWorker.autoCommit(batch))
			}
			//批次未满
			if len(batch.items) < 100 {
				continue
			}
			batch.commitTimer.Stop()
		case timeOutBatch = <-mergeWorker.timeOutChan:
			if mergeWorker.allBatch != timeOutBatch {
				continue
			}
		}
		mergeWorker.commitBatch(batch)
	}
}

func (mergeWorker *MergeWorker) autoCommit(batch *PushBatch) func() {
	return func() {
		mergeWorker.timeOutChan <- batch
	}
}

func (mergeWorker *MergeWorker) commitBatch(batch *PushBatch) {
	var (
		bizMsg *BizMessage
		buf    []byte
	)
	bizMsg = &BizMessage{
		Item: batch.items,
	}
	if buf, err = json.Marshal(*bizMsg); err != nil {
		return
	}
	G_ConnMgr.PushAll(buf)
}

func (mergeWorker *MergeWorker) PushAll(msg json.RawMessage) {
	var (
		pushContext *PushContext
	)
	pushContext = &PushContext{
		msg: &msg,
	}
	select {
	case mergeWorker.contextChan <- pushContext:
	default:

	}
}
