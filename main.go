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
	if conf.Config.Target == "" {
		return
	}
	for {
		online := utils.Check(conf.Config.Target, 5*time.Second)
		if online {
			message := "服务器在线"
			log.Println(message)
			updateDNSRecords(conf.Config.Mata.Main, message)
		} else {
			message := "服务器离线"
			log.Println(message)
			updateDNSRecords(conf.Config.Mata.Then, message)
		}
		if Once {
			return
		}
		time.Sleep(300 * time.Second)
	}
}

func updateDNSRecords(records []conf.DNSRecord, message string) {
	send := false
	for _, record := range records {
		ok, dns := utils.GetDnsRecoid(record.Name)
		if ok && dns.Content != record.Content {
			send = true
			log.Printf("修改解析【%s】\n", record.Name)
			utils.Dns(record, dns.ID)
		}
		time.Sleep(5 * time.Second)
	}
	if send && conf.Config.BotToken != "" && conf.Config.ChatID != "" {
		utils.SendMessage(message)
	}
}

func init() {
	once := flag.Bool("once", false, "Run once")
	flag.Parse()
	if *once {
		Once = true
	}
}
