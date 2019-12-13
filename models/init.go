package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"time"
)

type MysqlConfig struct {
	UserName string
	UserPwd  string
	Port     int
	Host     string
	DbName   string
}

type RedisConfig struct {
	Address     string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

type ConfigAll struct {
	MysqlConfig
	RedisConfig
}

var (
	DB        orm.Ormer
	G_Config  ConfigAll
	err       error
	appConf   config.Configer
	RedisPool *redis.Pool
)

func init() {
	orm.Debug = true
	orm.RegisterModel(new(EggUser), new(EggUserBalance))
	if err = InitAll(); err != nil {
		return
	}
	orm.RunSyncdb("default", false, true)
}

func InitAll() (err error) {
	G_Config, err = InitCfg()
	if err != nil {
		return
	}
	if DB, err = InitMysql(); err != nil {
		return
	}
	if RedisPool, err = InitRedis(); err != nil {
		return
	}
	return
}

func InitMysql() (Db orm.Ormer, err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", G_Config.MysqlConfig.UserName, G_Config.MysqlConfig.UserPwd, G_Config.MysqlConfig.Host, G_Config.MysqlConfig.Port, G_Config.MysqlConfig.DbName)
	fmt.Println("dns:::::", dns)
	if err = orm.RegisterDataBase("default", "mysql", dns); err != nil {
		return
	}
	Db = orm.NewOrm()
	return
}

func InitRedis() (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     G_Config.RedisConfig.MaxIdle,                    // 最初连接数量
		MaxActive:   G_Config.RedisConfig.MaxActive,                  // 最大连接数量 0表示按需创建
		IdleTimeout: time.Duration(G_Config.RedisConfig.IdleTimeout), // 连接关闭时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", G_Config.RedisConfig.Address)
		},
	}
	conn := pool.Get()
	defer func(conn redis.Conn) { conn.Close() }(conn)
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("redis connect failed")
		err = errors.New("redis connect failed")
		return
	}
	return
}

func InitCfg() (ConfigAll ConfigAll, err error) {
	appConf, err = config.NewConfig("ini", "conf/app.conf")
	//workPath, err := os.Getwd()
	//fmt.Println("workPath:::::::", workPath)
	if err != nil {
		err = errors.New("config parse error")
		return
	}
	if ConfigAll.MysqlConfig, err = GetMysqlConfig(); err != nil {
		return
	}
	if ConfigAll.RedisConfig, err = GetRedisConfig(); err != nil {
		return
	}
	//G_Config = ConfigAll
	return
}

func GetMysqlConfig() (MysqlConfig MysqlConfig, err error) {
	UserName := appConf.String("mysql::mysqluser")
	if len(UserName) == 0 {
		logs.Error("load config mysqluser failed")
		err = errors.New("load config mysqluser failed")
		return
	}
	MysqlConfig.UserName = UserName
	UserPwd := appConf.String("mysql::mysqlpass")
	if len(UserPwd) == 0 {
		logs.Error("load config mysqlpass failed")
		err = errors.New("load config mysqlpass failed")
		return
	}
	MysqlConfig.UserPwd = UserPwd
	Port, err := appConf.Int("mysql::mysqlport")
	if err != nil {
		logs.Error("load config mysqlport failed")
		err = errors.New("load config mysqlport failed")
		return
	}
	MysqlConfig.Port = Port
	Host := appConf.String("mysql::mysqlhost")
	if len(Host) == 0 {
		logs.Error("load config mysqlhost failed")
		err = errors.New("load config mysqlhost failed")
		return
	}
	MysqlConfig.Host = Host
	DbName := appConf.String("mysql::mysqldb")
	if len(DbName) == 0 {
		logs.Error("load config mysqldb failed")
		err = errors.New("load config mysqldb failed")
		return
	}
	MysqlConfig.DbName = DbName
	return
}

func GetRedisConfig() (RedisConfig RedisConfig, err error) {
	Address := appConf.String("redis::redis_addr")
	if len(Address) == 0 {
		logs.Error("load redis config failed")
		err = errors.New("load redis config redis_addr failed")
	}
	RedisConfig.Address = Address
	MaxIdle, err := appConf.Int("redis::redis_MaxIdle")
	if err != nil {
		logs.Error("load redis config failed")
		err = errors.New("load redis config redis_MaxIdle failed")
	}
	RedisConfig.MaxIdle = MaxIdle

	MaxActive, err := appConf.Int("redis::redis_MaxActive")
	if err != nil {
		logs.Error("load redis config failed")
		err = errors.New("load redis config redis_MaxActive failed")
	}
	RedisConfig.MaxActive = MaxActive
	IdleTimeout, err := appConf.Int("redis::redis_IdleTimeout")
	if err != nil {
		fmt.Println(err)
		logs.Error("load redis config redis_IdleTimeout failed")
		err = errors.New("load redis config redis_IdleTimeout failed")
	}
	RedisConfig.IdleTimeout = IdleTimeout

	return
}
