package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"csz.net/mata/conf"
	"csz.net/mata/utils"
)

var Once bool

func main() {
	for {
		for _, mata := range conf.Config.Mata {
			if mata.PS == "" {
				mata.PS = mata.Target
			}
			log.Println("开始检测" + mata.PS)
			send := false
			msg := "服务器在线"
			onlineCount := 0
			online := false
			for i := 0; i < 3; i++ {
				check, status := utils.Check(mata.Target, 5*time.Second)
				if !check {
					log.Println("检测失败")
					continue
				}
				if status {
					onlineCount++
				}
				time.Sleep(10 * time.Second)
			}
			if onlineCount >= 2 {
				online = true
			}
			if online {
				log.Println(msg)
				ok, dns := utils.GetDnsRecoid(mata.Main.Name, mata.Main.ZoneID)
				if ok && dns.Content != mata.Main.Content {
					send = true
					log.Printf("修改解析【%s】\n", mata.Main.Name)
					utils.Dns(mata.Main, dns.ID, mata.Main.ZoneID)
				} else {
					log.Printf("无需修改解析【%s】\n", mata.Main.Name)
				}
			} else {
				msg = "服务器离线"
				log.Println(msg)
				if mata.Then.ZoneID == "" {
					// 不填写then的zoneid则默认为main的zoneid
					mata.Then.ZoneID = mata.Main.ZoneID
				}
				ok, dns := utils.GetDnsRecoid(mata.Then.Name, mata.Then.ZoneID)
				if ok && dns.Content != mata.Then.Content {
					send = true
					log.Printf("修改解析【%s】\n", mata.Then.Name)
					utils.Dns(mata.Then, dns.ID, mata.Then.ZoneID)
				}
			}
			if send && conf.Config.BotToken != "" && conf.Config.ChatID != "" {
				msg = "【" + mata.PS + "】" + msg
				go utils.SendMessage("#MATA " + msg)
			}
		}
		if Once {
			return
		}
		time.Sleep(time.Duration(conf.Config.Corn) * time.Second)
	}
}
func init() {
	once := flag.Bool("once", false, "Run once")
	configPath := flag.String("config", "mata.json", "Config file path")
	flag.Parse()
	if *once {
		Once = true
	}
	configInit(*configPath)
}
func configInit(path string) {
	// 读取配置文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(path + "配置文件不存在")
		os.Exit(0)
	}
	defer file.Close()

	// 解码 JSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf.Config)
	if err != nil {
		fmt.Println(path + "配置文件错误")
		os.Exit(0)
	}
	conf.Config.TgApiUrl = strings.TrimRight(conf.Config.TgApiUrl, "/")
}
