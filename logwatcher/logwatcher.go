package main

import (
	"bytes"
	"encoding/json"
	"github.com/hpcloud/tail"
	"github.com/londonhackspace/acnode-dashboard/api"
	"net/http"
	"os"
	"regexp"
	"time"
)

const regexmatcher string = "[0-9\\.]+\\s-\\s(?:.*)\\s\\[([0-9\\/A-Za-z\\:\\+\\s]+)\\]\\s\\\"GET\\s\\/([0-9a-zA-Z]+)\\/[a-zA-Z\\/]+\\sHTTP\\/[0-9\\.]+\\\"\\s[0-9]{3}\\s[0-9]+\\s\"[0-9A-Za-z\\/\\:-]+\"\\s\"ACNode rev ([a-zA-Z0-9\\-]+)\""

func main() {
	cfg := GetConfig(os.Args[1])

	t, err := tail.TailFile(cfg.LogFile, tail.Config{Follow: true, ReOpen: true, MustExist: true})
	if err != nil {
		panic(err)
	}

	matcher := regexp.MustCompile(regexmatcher)
	for line := range t.Lines {
		res := matcher.FindStringSubmatch(line.Text)
		if len(res) == 4 {
			version := res[3]
			nodeId := res[2]
			date,_ := time.Parse("2/Jan/2006:15:04:05 -0700", res[1])
			var b api.SetStatusBody
			b.Version = version
			b.Timestamp = date.Unix()

			body,_ := json.Marshal(b)

			client := &http.Client{}
			bodyreader := bytes.NewReader(body)
			req, _ := http.NewRequest("POST", cfg.ACNodeDashApiUrl + "nodes/setStatus/"+nodeId, bodyreader)
			req.Header.Add("API-KEY", cfg.ACNodeDashApiKey)
			client.Do(req)
		}
	}
}
