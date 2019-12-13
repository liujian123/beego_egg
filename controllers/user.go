package controllers

import (
	"apiproject/models"
	"apiproject/utils"
	"fmt"
	"strings"
	"time"
)

// Operations about Users
type UserController struct {
	BaseController
}

var (
	UserModel *models.EggUser = models.NewUserModel()
)

func (this *UserController) URLMapping() {
	this.Mapping("Login", this.Login)
	this.Mapping("Info", this.Info)
	this.Mapping("Collect", this.Collect)
	this.Mapping("Withdraw", this.Withdraw)
	this.Mapping("Balance", this.Balance)
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (this *UserController) Login() {
	username := this.GetString("username")
	password := this.GetString("password")
	user, err := UserModel.GetUserByNameAndPwd(username, password)
	if err == nil {
		et := utils.EasyToken{
			Username: user.UserName,
			Uid:      user.UserId,
			Expires:  time.Now().Unix() + 3600*100,
		}
		token, err := et.GetToken()
		if token == "" || err != nil {
			this.Data["json"] = &ErrResponse{16001, fmt.Sprintf("%s", err)}
		} else {
			this.Data["json"] = &Response{0, "succ", user}
			this.Ctx.Output.Header("Authorization", token)
		}
	} else {
		this.Data["json"] = &ErrResponse{16002, fmt.Sprintf("%s", err)}
	}
	this.ServeJSON()
}

func (this *UserController) Info() {
	//authorization := this.Ctx.Input.Header("authorization")
	uid := this.Ctx.Input.GetData("uid").(int64)
	user, err := UserModel.GetUserInfo(uid)
	if err != nil {
		this.Data["json"] = &ErrResponse{16001, fmt.Sprintf("%s", err)}
	} else {
		this.Data["json"] = &Response{0, "succ", user}
	}
	this.ServeJSON()
}

func (this *UserController) Collect() {
	keys := strings.Split(this.GetString("keys"), ",")
	if len(keys) == 0 {
		this.Data["json"] = &ErrResponse{16001, "参数有误!"}
		this.ServeJSON()
		return
	}
	uid := this.Ctx.Input.GetData("uid").(int64)
	//RequestBody := this.Ctx.Input.RequestBody
	result, err := UserModel.CollectEggs(uid, keys)
	if err != nil {
		this.Data["json"] = &ErrResponse{16001, fmt.Sprintf("%s", err)}
	} else {
		this.Data["json"] = &Response{0, "succ", result}
	}
	this.ServeJSON()
}

func (this *UserController) Withdraw() {
	weight, err := this.GetInt("weight")
	if err != nil || weight != 50 {
		this.Data["json"] = &ErrResponse{16001, "参数有误"}
		this.ServeJSON()
		return
	}
	uid := this.Ctx.Input.GetData("uid").(int64)
	result, err := UserModel.WithdrawApply(uid, weight)
	if err != nil {
		this.Data["json"] = &ErrResponse{16001, fmt.Sprintf("%s", err)}
	} else {
		this.Data["json"] = &Response{0, "succ", result}
	}
	this.ServeJSON()
}

func (this *UserController) Balance() {
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
				this.Data["json"] = &ErrResponse{16001, "Error: invalid query key/value pair"}
				this.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	balance, err := UserModel.GetBalanceList(query, fields, sortby, order, offset, limit)
	if err != nil {
		this.Data["json"] = &ErrResponse{16002, "获取数据有误"}
	} else {
		this.Data["json"] = &Response{0, "succ", balance}

	}
	this.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout/:uid [get,post]
func (this *UserController) Logout() {
	this.Data["json"] = "logout success"
	this.ServeJSON()
}
