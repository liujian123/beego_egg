package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func RegisterDB() {
	user := beego.AppConfig.String("mysqluser")
	pwd := beego.AppConfig.String("mysqlpass")
	url := beego.AppConfig.String("mysqlurls")
	db := beego.AppConfig.String("mysqldb")
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pwd, url, db)
	fmt.Println("DB_URL::::::::::", dns)
	fmt.Println("DROracle::::::::::", orm.DROracle)
	fmt.Println("DRTiDB::::::::::", orm.DRTiDB)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(Topic))
	orm.RegisterDataBase("default", "mysql", dns)
	// create table
	orm.RunSyncdb("default", false, true)
}
