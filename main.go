package main

import (
	"flag"
	"log"
	"time"

	"csz.net/mata/conf"
	"csz.net/mata/utils"
)

var Sus bool

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
		if Sus {
			return
		}
		time.Sleep(10 * time.Second)
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
	sus := flag.Bool("sus", false, "Disable Sustain")
	flag.Parse()
	if *sus {
		Sus = true
	}
}
