package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

func getConfigFile() []byte {
	path, err := homedir.Expand("~/lottery_conf.yaml")
	if err != nil {
		panic("get homedir failed")
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("read file failed")
	}
	return file
}

var m map[interface{}]interface{} = nil

// GetConfig get config in lottery_conf.yaml
func GetConfig(namespace string) map[interface{}]interface{} {
	if m == nil {
		m = make(map[interface{}]interface{})
		err := yaml.Unmarshal(getConfigFile(), &m)
		if err != nil {
			panic(err)
		}
	}
	return m[namespace].(map[interface{}]interface{})
}

func getMysqlSource() string {
	config := GetConfig("mysql")
	if config["password"] == nil {
		return fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config["user"],
			config["host"],
			config["port"],
			config["database"],
		)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["user"],
		config["password"],
		config["host"],
		config["port"],
		config["database"],
	)
}

func getRedisOption() *redis.Options {
	config := GetConfig("redis")
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config["host"], config["port"]),
		Password: config["password"].(string),
		DB:       config["database"].(int),
	}
}

func getKafkaAddr() (addr []string) {
	config := GetConfig("kafka")
	ai := config["addr"].([]interface{})
	for _, i := range ai {
		addr = append(addr, i.(string))
	}
	return
}
