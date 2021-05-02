package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type LogWatcherConfig struct {
	LogFile string `json:"LogFile"`
	Follow bool `json:"follow"`
	ACNodeDashApiUrl string `json:"acnodedash-apiurl"`
	ACNodeDashApiKey string `json:"acnodedash-apikey"`
}

func GetConfig(filename string) LogWatcherConfig {
	f,err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data,err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	ret := LogWatcherConfig{
		Follow: true,
	}

	json.Unmarshal(data, &ret)

	return ret
}