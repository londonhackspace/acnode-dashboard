package main

import (
	"io/ioutil"
)

var version *string = nil

func getVersion() string {
	if version == nil {
		stringVer := "Unknown"
		if content, err := ioutil.ReadFile("version"); err == nil {
			stringVer = string(content)
		}
		version = &stringVer
	}

	return *version
}