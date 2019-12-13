package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

// topic info
type Topic struct {
	Id         int    `json:"id" orm:"column(id);auto"`
	Uid        int64  `json:"uid" orm:"column(uid)"`
	Title      string `json:"title" orm:"column(title);size(255)"`
	Content    string `json:"content" orm:"column(content);type(text)"`
	Summary    string `json:"summary" orm:"column(summary);size(200)"`
	Attachment string `json:"url" orm:"column(attachment);size(255)"`
	//Category        *Category `json:"cate" orm:"rel(fk);on_delete(do_nothing)"`
	//Labels          []*Label  `json:"labels" orm:"rel(m2m)"`
	Created         time.Time `orm:"auto_now_add;column(created);type(datetime)"`
	Updated         time.Time `json:"-" orm:"auto_now;column(updated);type(datetime)"`
	Deleted         time.Time `json:"-" orm:"auto_now;column(deleted);type(datetime)"`
	Views           int64     `json:"views" orm:"column(views)"`
	Author          string    `json:"author" orm:"column(author);size(255)"`
	ReplyTime       time.Time `json:"-" orm:"column(reply_time);type(datetime);null"`
	ReplyCount      int64     `json:"reply_count" orm:"column(reply_count)"`
	ReplyLastUserId int64     `json:"reply_last_user_id" orm:"column(reply_last_user_id)"`
}

//自定义表名
func (t *Topic) TableName() string {
	return "topic"
}

func AddTopic(m *Topic) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return id, err
	}
	return id, err
}

func GetAllTopics(query map[string]string, fields []string, sortby []string, order []string, offset, limit int64) ([]*Topic, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Topic))
	for k, v := range query {
		qs = qs.Filter(k, v)
	}
	var sortFields []string
	if len(sortby) == len(order) {
		orderby := ""
		for i, v := range sortby {
			if order[i] == "desc" {
				orderby = "-" + v
			} else {
				orderby = v
			}
			sortFields = append(sortFields, orderby)
		}
	}
	fmt.Println("sortFields::::", sortFields)
	qs = qs.OrderBy(sortFields...)
	var err error
	var topics []*Topic
	if _, err = qs.Limit(limit, offset).All(&topics, fields...); err == nil {
		//for _, topic := range topics {
		//	o.Read(topic)
		//	if err != nil {
		//		return nil, err
		//	}
		//}
		fmt.Println("topics::::::::::====================", &topics)
		return topics, nil
	}
	return nil, err
}

func UpdateTopicById(m *Topic) (err error) {
	o := orm.NewOrm()
	var num int64
	if num, err = o.Update(m); err == nil {
		fmt.Println("num:::::", num)
		fmt.Println("m:::::", m)
	}
	return
}

func DeleteTopicById(id int) (err error) {
	o := orm.NewOrm()
	v := &Topic{Id: id}
	if err = o.Read(v); err == nil {
		fmt.Println("v:::::", v)
		var num int64
		if num, err = o.Delete(&Topic{Id: id}); err == nil {
			fmt.Println("num:::::", num)
		}
	}
	return
}

func GetTopicByID(id int) (v *Topic, err error) {
	o := orm.NewOrm()
	v = &Topic{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}
