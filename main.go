package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"

	"gopkg.in/telegram-bot-api.v4"
)

type Config struct {
	Token       string `json:"token"`
	ChannelName string `json:"channelName"`
}

var configPath = "config.json"
var bot *tgbotapi.BotAPI

const ojiPage = "http://olimpiada.info/oji2017/index.php?cid=rezultate"
const sleepSeconds = 1

func (c *Config) readConfig(fp string) error {
	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, c)
	if err != nil {
		return err
	}

	return nil
}

func fetchWebPage(url string) string {
	resp, err := http.Get(url)
	if err == nil {
		buf, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			return string(buf)
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func diffEqual(diff []diffmatchpatch.Diff) bool {
	for _, d := range diff {
		if d.Type != diffmatchpatch.DiffEqual {
			return false
		}
	}
	return true
}

func main() {
	var conf Config
	err := conf.readConfig(configPath)
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		panic(err)
	}
	for true {
		page1 := fetchWebPage(ojiPage)
		time.Sleep(sleepSeconds * time.Second)
		page2 := fetchWebPage(ojiPage)
		dmp := diffmatchpatch.New()
		diff := dmp.DiffMain(page1, page2, false)
		if !diffEqual(diff) {
			log.Println("Sending message")
			msg := tgbotapi.NewMessageToChannel(conf.ChannelName, "S-a actualizat pagina!")
			bot.Send(msg)
		} else {
			log.Println("Pages equal")
		}
		log.Println("Zzz...")
		time.Sleep(sleepSeconds * time.Second)
	}
}
