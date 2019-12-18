package models

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	msgType int
	msgData []byte
}

type WsConnection struct {
	WsSocket *websocket.Conn
	ConnId   uint64
	OutChan  chan *WsMessage
	//CloseChan
}

func InitWsConnection(WsSocket *websocket.Conn, connId uint64) *WsConnection {
	WsConn := &WsConnection{
		WsSocket: WsSocket,
		ConnId:   connId,
		OutChan:  make(chan *WsMessage, 100),
	}
	go WsConn.ReadLoop()
	go WsConn.WriteLoop()
	return WsConn
}

func (WsConnection *WsConnection) ReadLoop() {
	var (
		msgType int
		data    []byte
		err     error
	)
	for {
		if msgType, data, err = WsConnection.WsSocket.ReadMessage(); err != nil {
			fmt.Println(msgType, string(data), err)
		}
	}
}

func (WsConnection *WsConnection) WriteLoop() {
	var (
		WsMsg *WsMessage
		err   error
	)
	for {
		select {
		case WsMsg = <-WsConnection.OutChan:
			if err = WsConnection.WsSocket.WriteMessage(WsMsg.msgType, WsMsg.msgData); err != nil {
				goto ERR
			}
		}
	}
ERR:
	fmt.Println("WsSocket Write failed")
}

func (WsConnection *WsConnection) SendMessage(pushJob *PushJob) {
	var (
		WsMsg *WsMessage
	)
	fmt.Println("SendMessage:SendMessage")
	WsMsg = &WsMessage{
		msgType: pushJob.MsgType,
		msgData: pushJob.MsgData,
	}
	select {
	case WsConnection.OutChan <- WsMsg:
	}
}

func (WsConnection *WsConnection) WSHandle(ConnId uint64) {
	G_ConnMgr.AddConn(WsConnection)
}
