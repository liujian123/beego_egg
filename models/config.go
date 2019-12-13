package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	G_config map[string]map[string]map[string]SourceConf
)

type Pos struct {
	SCREEN string `json:"SCREEN"`
	FEED   string `json:"FEED"`
}

type SourceConf struct {
	Appid   string `json:"appid"`
	Appname string `json:"appname"`
	Pkgname string `json:"pkgname"`
	Appver  string `json:"appver"`
	Pos     Pos    `json:"pos"`
}

func InitConfig() (err error) {
	var (
		content  []byte
		conf     map[string]map[string]map[string]SourceConf
		workPath string
		confPath string
	)

	if workPath, err = os.Getwd(); err != nil {
		return
	}
	confPath = filepath.Join(workPath, "conf", "config.json")
	if content, err = ioutil.ReadFile(confPath); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = conf
	return
}
