package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
}

//Response 结构体
type Response struct {
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

//Response 结构体
type ErrResponse struct {
	Errcode int         `json:"errcode"`
	Errmsg  interface{} `json:"errmsg"`
}
