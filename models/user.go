package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"math"
	"strconv"
	"strings"
	"time"
)

type EggUser struct {
	UserId     int64     `json:"uid" orm:"pk;auto"`
	AgentUid   int64     `json:"agent_uid" orm:"default(0)"`
	EggsActive float64   `json:"eggs_active" orm:"default(0)"`
	EggsTotal  float64   `json:"eggs_total" orm:"default(0)"`
	UserName   string    `json:"username" orm:"size(60);default('')" form:"UserName"  valid:"Required"`
	UserPwd    string    `json:"pwd" orm:"size(60);default('')" form:"UserPwd"  valid:"Required"`
	Avatar     string    `json:"avatar" orm:"size(60);default('')" form:"Avatar"  valid:"Required"`
	Chickens   int       `json:"chickens" orm:"default(0)" form:"Chickens"`
	UserEmail  string    `json:"email" orm:"size(60);default('')" form:"UserEmail"  valid:"Email"`
	UserMobile uint64    `json:"mobile" orm:"size(11);default(0)" form:"UserMobile" valid:"Mobile"`
	Status     uint      `json:"status" orm:"default(1)" form:"Status" valid:"Range(1,3)"`
	IsLogin    uint      `json:"is_login" orm:"default(1)" form:"IsLogin" valid:"Range(1,3)"`
	CreateTime time.Time `json:"create_time" orm:"size(10);default(0)" form:"CreateTime" valid:"Min(0)"`
}

type EggUserBalance struct {
	Id         int64   `json:"id" orm:"pk;auto"`
	Uid        int64   `json:"uid" orm:"default(0)"`
	Type       int     `json:"type" orm:"default(0)"`
	BizType    string  `json:"biz_type" orm:"size(60);default('')"`
	Ymd        int     `json:"ymd" orm:"default(0)"`
	Amount     float64 `json:"amount" orm:"default(0)"`
	Balance    float64 `json:"balance" orm:"default(0)"`
	Extra      string  `json:"email" orm:"size(150);default('')"`
	Status     uint    `json:"status" orm:"default(1)" form:"Status"`
	CreateTime int64   `json:"create_time" orm:"default(0)"`
	UpdateTime int64   `json:"update_time" orm:"default(0)"`
}

type Eggs struct {
	Chickens int     `json:"chickens"`
	Eggs     float64 `json:"eggs"`
	S_Eggs   float64 `json:"s_eggs"`
	Key      string  `json:"key"`
}

type ResponseUserInfo struct {
	UserId     int64    `json:"uid"`
	UserName   string   `json:"username"`
	UserMobile uint64   `json:"mobile"`
	InviteCode int64    `json:"invite_code"`
	AgentUid   int64    `json:"agent_uid"`
	Avatar     string   `json:"avatar"`
	Chickens   int      `json:"chickens"`
	SystemTime int64    `json:"system_time"`
	EggsActive float64  `json:"eggs_active"`
	EggsTotal  float64  `json:"eggs_total"`
	Flag       int      `json:"flag"`
	SysMsg     []string `json:"sys_msg"`
	Eggs       []*Eggs  `json:"eggs"`
}

type Msg struct {
	Uid        int64
	UserName   string
	CreateTime string
}

var (
	MsgList = make([]string, 0, 20)
	MsgChan = make(chan *Msg, 20)
)

func NewUserModel() *EggUser {
	return &EggUser{}
}

func init() {
	go LoadSetting()
	go SyncMsg()
}

func (this *EggUser) GetUserByNameAndPwd(username, password string) (user *EggUser, err error) {
	user = &EggUser{
		UserName: username,
		UserPwd:  password,
	}
	if err = DB.Read(user, "UserName", "UserPwd"); err != nil {
		if err == orm.ErrNoRows {
			err = errors.New("用户名或密码错误!")
			return
		}
		err = errors.New(fmt.Sprintf("GetUserByNameAndPwd:err:%v,username:%s,password:%s", err, username, password))
		logs.Warn(err)
		return
	}
	return
}

func (this *EggUser) GetUserInfo(uid int64) (respUserInfo *ResponseUserInfo, err error) {
	user, err := Find(uid)
	respUserInfo = &ResponseUserInfo{
		UserId:     user.UserId,
		UserName:   user.UserName,
		UserMobile: user.UserMobile,
		InviteCode: user.UserId,
		AgentUid:   user.AgentUid,
		Avatar:     user.Avatar,
		Chickens:   user.Chickens,
		SystemTime: time.Now().Unix(),
		EggsActive: user.EggsActive,
		EggsTotal:  user.EggsTotal,
		Flag:       0,
		SysMsg:     MsgList,
		//Eggs:       make([]*Eggs, 0, 20),
	}
	if respUserInfo.Eggs, err = getUserEggs(uid); err != nil {
		return
	}
	return
}

func getUserEggs(uid int64) (eggs []*Eggs, err error) {
	strUid := strconv.Itoa(int(uid))
	conn := RedisPool.Get()
	key := "eggs:lay:" + strUid
	result, err := redis.StringMap(conn.Do("hgetall", key))
	if err != nil {
		logs.Error("hgetall err")
		err = errors.New("hgetall err")
		conn.Close()
		return
	}
	for k, v := range result {
		if v == "" {
			continue
		}
		var (
			eggsList []interface{}
			egg      Eggs
		)
		err = json.Unmarshal([]byte(v), &eggsList)
		if err != nil {
			return
		}
		egg.Chickens = int(eggsList[1].(float64))
		egg.Eggs = eggsList[2].(float64)
		egg.S_Eggs = eggsList[3].(float64)
		egg.Key = k
		eggs = append(eggs, &egg)
	}

	if eggs == nil {
		eggs = make([]*Eggs, 0)
	}
	conn.Close()
	return
}

func Find(uid int64) (user *EggUser, err error) {

	//方法1
	/*err = DB.Raw("SELECT * FROM `egg_user` WHERE `user_id` = ?", uid).QueryRow(&user)
	fmt.Println("user:::::", user)*/

	//方法二
	user = &EggUser{
		UserId: uid,
	}
	if err = DB.Read(user, "UserId"); err != nil {
		if err == orm.ErrNoRows {
			err = errors.New("用户名或密码错误!")
			return
		}
		err = errors.New(fmt.Sprintf("GetUserByNameAndPwd:err:%v,username:%s,password:%s", err, user.UserName, user.UserPwd))
		logs.Warn(err)
		return
	}

	return
}

func (this *EggUser) CollectEggs(uid int64, keys []string) (result map[string]float64, err error) {
	strUid := strconv.Itoa(int(uid))
	user, err := Find(uid)
	conn := RedisPool.Get()
	key := "eggs:lay:" + strUid
	eggs, err := redis.StringMap(conn.Do("hgetall", key))

	if err != nil {
		conn.Close()
		return
	}
	fmt.Println("keys", keys)
	var (
		weight   float64 = 0
		eggsData         = make(map[string]string)
	)
	for _, v := range keys {
		var (
			eggsList []interface{}
		)
		if vv, ok := eggs[v]; ok && vv != "" {
			err = json.Unmarshal([]byte(vv), &eggsList)
			if err != nil {
				conn.Close()
				logs.Error("err:::", err)
				return
			}
			eggsData[v] = ""
			weight += eggsList[2].(float64) - eggsList[3].(float64)
		}
	}
	_, err = conn.Do("hmset", redis.Args{}.Add(key).AddFlat(eggsData)...)
	if err != nil {
		conn.Close()
		logs.Error("err:::", err)
		err = errors.New("hmset eggsData failed")
		return
	}
	fmt.Println("weight", weight)
	result = make(map[string]float64, 1)
	if weight == 0 {
		conn.Close()
		logs.Error("err:::", err)
		result["eggs_active"] = user.EggsActive
		return
	}
	weight = math.Round(weight*100000) / 100000
	_, err = DB.Raw("UPDATE egg_user SET eggs_active = eggs_active +?, eggs_total=eggs_total+? WHERE user_id = ?", weight, weight, uid).Exec()
	if err != nil {
		conn.Close()
		return
	}
	result["eggs_active"] = user.EggsActive + weight
	conn.Close()
	return
}

func (this *EggUser) WithdrawApply(uid int64, weight int) (result map[string]float64, err error) {
	var (
		updateErr error
		insertErr error
	)
	//todo 事务
	err = DB.Begin()
	user, err := Find(uid)
	_, updateErr = DB.Raw("UPDATE egg_user SET eggs_active = eggs_active - ? WHERE user_id = ?", weight, uid).Exec()
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	ymd, _ := strconv.Atoi(year + month + day)
	amount := float64(weight) * 0.16
	UserBalance := &EggUserBalance{
		Uid:        uid,
		Type:       1,
		BizType:    "withdraw",
		Ymd:        ymd,
		Amount:     amount,
		Balance:    0,
		Extra:      "",
		Status:     1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	_, insertErr = DB.Insert(UserBalance)
	if updateErr != nil || insertErr != nil {
		DB.Rollback()
		err = errors.New("提现失败")
	} else {
		DB.Commit()
		result = make(map[string]float64, 1)
		result["eggs_active"] = user.EggsActive - float64(weight)
		msg := &Msg{
			Uid:        uid,
			UserName:   user.UserName,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		MsgChan <- msg
	}
	return
}

func (this *EggUser) GetBalanceList(query map[string]string, fields []string, sortby []string, order []string, offset, limit int64) (balance []*EggUserBalance, err error) {
	qs := DB.QueryTable(new(EggUserBalance))
	for k, v := range query {
		qs = qs.Filter(k, v)
	}
	var (
		sortField []string
	)
	if len(sortby) == len(order) {
		for i := range sortby {
			if strings.ToLower(order[i]) == "desc" {
				sortField = append(sortField, "-"+sortby[i])
			} else {
				sortField = append(sortField, sortby[i])
			}
		}
	}
	qs.OrderBy(sortField...)
	_, err = qs.Limit(limit, offset).All(&balance, fields...)
	if err != nil {
		return
	}
	return
}

func SyncMsg() {
	for {
		msg, ok := <-MsgChan
		if !ok {
			time.Sleep(time.Millisecond * 100) // 停止100毫秒
			continue
		}
		conn := RedisPool.Get()
		key := "eggs:msg"
		msgs_num, err := redis.Int(conn.Do("LLEN", key))
		if err != nil {
			conn.Close()
			continue
		}
		strMsg := `"` + "恭喜" + msg.UserName + "在" + msg.CreateTime + " 出售鸡蛋" + `"`
		if msgs_num >= 20 {
			fmt.Println("阻塞等待返回::::::", msgs_num)
			//todo 阻塞等待返回，设置超时时间直接返回
			reply, err := conn.Do("BRPOP", "eggs:msg", 1)
			data, err := redis.Strings(reply, err)
			fmt.Println(data)
			if err == redis.ErrNil {
				logs.Warn("[SyncMsg] 未pop出数据")
			}
		}
		_, err = conn.Do("LPUSH", key, strMsg)
		MsgList = append(MsgList, strMsg)
		conn.Close()
		//推送给所有的在线链接（客户端）
		G_merger.broadcastWorker.PushAll(json.RawMessage(strMsg))
	}
}

func LoadSetting() (err error) {
	conn := RedisPool.Get()
	key := "eggs:msg"
	MsgList, err = redis.Strings(conn.Do("LRANGE", key, 0, -1))
	conn.Close()
	if err != nil {
		logs.Error("LoadSetting failed", err)
		err = errors.New("LoadSetting failed")

	}
	return
}
