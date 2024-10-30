package main

import (
	"flag"
	"log"
	"time"

	"csz.net/mata/conf"
	"csz.net/mata/utils"
)

var Once bool

func main() {
	for {
		for _, mata := range conf.Config.Mata {
			log.Println("开始检测" + mata.Target)
			send := false
			msg := "服务器在线"
			onlineCount := 0
			online := false
			for i := 0; i < 3; i++ {
				if utils.Check(mata.Target, 5*time.Second) {
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
				}
			} else {
				msg = "服务器离线"
				log.Println(msg)
				ok, dns := utils.GetDnsRecoid(mata.Then.Name, mata.Then.ZoneID)
				if ok && dns.Content != mata.Then.Content {
					send = true
					log.Printf("修改解析【%s】\n", mata.Then.Name)
					utils.Dns(mata.Then, dns.ID, mata.Then.ZoneID)
				}
			}
			if send && conf.Config.BotToken != "" && conf.Config.ChatID != "" {
				msg = "【" + mata.Target + "】" + msg
				go utils.SendMessage(msg)
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
	flag.Parse()
	if *once {
		Once = true
	}
}
