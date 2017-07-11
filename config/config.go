package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

//Table keep config data
type Table map[string]interface{}

//Get : get config object
func (table Table) Get(key string) interface{} {
	if table.Contain(key) {
		return table[key]
	}
	return nil
}

//Get : get config string
func (table Table) GetString(key string) (string, error) {
	if table.Contain(key) {
		if val, ok := table[key].(string); ok {
			return val, nil
		}
		return "", errors.New("Value is not string")
	}
	return "", errors.New("Key Not Found")
}

//Get : get config bool
func (table Table) GetBoolean(key string) (bool, error) {
	if table.Contain(key) {
		if val, ok := table[key].(bool); ok {
			return val, nil
		}
		return false, errors.New("Value is not Boolean")
	}
	return false, errors.New("Key Not Found")
}

//Get : get int value
func (table Table) GetInt(key string) (int, error) {
	if table.Contain(key) {
		b, err := json.Marshal(table[key])
		if err != nil {
			return 0, err
		}
		var num int
		json.Unmarshal(b, &num)
		return num, nil
	}
	return 0, errors.New("Key Not Found")
}

//Add : add object config
func (table Table) Add(key string, obj interface{}) {
	table[key] = obj
}

//Contain : check key is contain in Table
func (table Table) Contain(key string) bool {
	if _, ok := table[key]; ok {
		return true
	}
	return false
}

//LoadConfig : load config from file
func LoadConfig(configFilePath string) (config *Table) {
	log.Println("Load config = ", configFilePath)
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic("cannot load config")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic("config invalid format")
	}

	return
}
