package utils

import (
	"fmt"
	"gin-project/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func InitReadConfigFile() {
	pwd, _ := os.Getwd()
	filePath := pwd + "/config.yaml"
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	conf := new(config.Server)
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		fmt.Println("解码错误: ", err)
		//错误处理
		return
	}
	fmt.Println(conf.Mysql)
}
