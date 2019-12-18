package controllers

import (
	"beego_egg/models"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
)

type WsSocketController struct {
	BaseController
}

var (
	WsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	CurConnId uint64 = uint64(time.Now().Unix())
)

func (this *WsSocketController) URLMapping() {
	this.Mapping("HandleConnect", this.HandleConnect)
}

func (this *WsSocketController) HandleConnect() {

	var (
		WsSocket *websocket.Conn
		WsConn   *models.WsConnection
		err      error
	)
	WsSocket, err = WsUpgrader.Upgrade(this.Ctx.Output.Context.ResponseWriter, this.Ctx.Output.Context.Request, nil)
	if err != nil {
		this.Data["json"] = &ErrResponse{16002, fmt.Sprintf("%s", err)}
		this.ServeJSON()
		return
	}

	CurConnId = atomic.AddUint64(&CurConnId, 1)
	WsConn = models.InitWsConnection(WsSocket, CurConnId)
	WsConn.WSHandle(CurConnId)
}
