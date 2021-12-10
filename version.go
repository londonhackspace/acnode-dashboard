package main

import (
	"io/ioutil"
	"strings"
)

var version *string = nil

func getVersion() string {
	if version == nil {
		stringVer := "Unknown"
		if content, err := ioutil.ReadFile("version"); err == nil {
			stringVer = strings.TrimSuffix(string(content), "\n")
		}
		version = &stringVer
	}

	return *version
}
