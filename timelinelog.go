package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Consumer_Key        string
	Consumer_Secret     string
	Access_Token        string
	Access_Token_Secret string
}

func Parse(filename string) (Config, error) {
	var c Config
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error" + err.Error())
		return c, err
	}
	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		fmt.Println("error" + err.Error())
		return c, err
	}
	return c, nil

}

func main() {
	conf, err := Parse("config.json")
	if err != nil {
		fmt.Println("error" + err.Error())
		return
	}
	anaconda.SetConsumerKey(conf.Consumer_Key)
	anaconda.SetConsumerSecret(conf.Consumer_Secret)

	api := anaconda.NewTwitterApi(conf.Access_Token, conf.Access_Token_Secret)
	api.SetLogger(anaconda.BasicLogger)

	v := url.Values{}
	stream := api.UserStream(v)
	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				t := time.Now()
				dir := fmt.Sprintf("%d/%d/%d/", t.Year(), int(t.Month()), t.Day())
				os.MkdirAll("./tweets/"+dir, 0755)
				b, err := json.Marshal(status)
				if err != nil {
					fmt.Println(err)
					continue
				}
				ioutil.WriteFile("./tweets/"+dir+strconv.FormatInt(status.Id, 10)+".json", b, os.ModePerm)
			default:
			}

		}

	}

}
