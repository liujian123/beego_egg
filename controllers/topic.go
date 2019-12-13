package controllers

import (
	"apiproject/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) URLMapping() {
	this.Mapping("GetAll", this.GetAll)
}

// @Title Create
// @Description Create Topics
// @Success 200 {object} models.Topic
// @router /create [post]
func (this *TopicController) Create() {
	var v models.Topic
	Result := models.NewCommResult()
	fmt.Println("this.Ctx.Input.RequestBody:::::::::", string(this.Ctx.Input.RequestBody))
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddTopic(&v); err != nil {
			Result.Msg = err.Error()
		} else {
			this.Ctx.Output.SetStatus(201)
			Result.Msg = "OK"
		}
	} else {
		Result.Msg = err.Error()
	}
	this.Data["json"] = Result
	this.ServeJSON()
}

// @Title GetAll
// @Description get all Topics
// @Success 200 {object} models.Topic
// @router / [get]
func (this *TopicController) GetAll() {

	//redis 测试
	//models.SetKey("nickname", "CheerYoung")
	var fields []string
	var sortby []string
	var order []string
	var offset int64
	var limit int64
	var query map[string]string = make(map[string]string)
	if v := this.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	if v, err := this.GetInt64("offset"); err == nil {
		offset = v
	}
	if v, err := this.GetInt64("limit"); err == nil {
		limit = v
	}
	if v := this.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	if v := this.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := this.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				this.Data["json"] = errors.New("Error: invalid query key/value pair")
				this.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	if topics, err := models.GetAllTopics(query, fields, sortby, order, offset, limit); err != nil {
		this.Data["json"] = err.Error()
	} else {
		this.Data["json"] = topics

	}

	this.ServeJSON()
}

// @Title Update
// @Description Update  Topics
// @Success 200 {object} models.Topic
// @router /update [post]
func (this *TopicController) Update() {
	v := models.Topic{}
	Result := models.NewCommResult()
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTopicById(&v); err == nil {
			Result.Msg = "ok"
		} else {
			Result.Msg = err.Error()
		}
	}

	this.Data["json"] = Result
	this.ServeJSON()
}

// @Title Delete
// @Description Update  Topics
// @Success 200 {object} models.Topic
// @router /delete [post]
func (this *TopicController) Delete() {
	v := models.Topic{}
	Result := models.NewCommResult()
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.DeleteTopicById(v.Id); err == nil {
			Result.Msg = "ok"
		} else {
			Result.Msg = err.Error()
		}
	}
	this.Data["json"] = Result
	this.ServeJSON()
}

// @Title GetOne
// @Description Get  Topic
// @Success 200 {object} models.Topic
// @router /:id [get]
func (this *TopicController) GetOne() {

	oldStr := this.Ctx.Input.Param(":id")
	fmt.Println("oldStr::::", oldStr)
	id, _ := strconv.Atoi(oldStr)
	if topic, err := models.GetTopicByID(id); err == nil {
		this.Data["json"] = topic
	} else {
		this.Data["json"] = err.Error()
	}

	this.ServeJSON()
}
