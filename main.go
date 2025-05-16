package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"csz.net/mata/conf"
	"csz.net/mata/utils"
)

var Once bool
var webOnce sync.Once // 确保 Web 只启动一次

func main() {
	for {
		var wg sync.WaitGroup

		for _, mata := range conf.Config.Mata {
			wg.Add(1)
			go func(mata conf.Mata) {
				defer wg.Done()

				if mata.PS == "" {
					mata.PS = mata.Target
				}
				log.Println("开始检测 " + mata.PS)

				send := false
				msg := "服务器在线"
				onlineCount := 0
				online := false

				if !strings.ContainsAny(mata.Target, ".") {
					// 启动 Web，只执行一次
					webOnce.Do(func() {
						go utils.Web()
						log.Println("被动检测端口: " + conf.WebPort)
					})

					if value, exists := conf.Array[mata.Target]; !exists {
						conf.Array[mata.Target] = int(time.Now().Unix()) + int(conf.Config.Corn)
						online = true
						log.Printf("被动检测, %s\n", mata.Target)
					} else {
						if value+int(conf.Config.Corn) >= int(time.Now().Unix()) {
							online = true
						}
					}
				} else {
					for i := 0; i < 3; i++ {
						check, status := utils.Check(mata.Target, 3*time.Second)
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
						mata.Then.ZoneID = mata.Main.ZoneID
					}
					ok, dns := utils.GetDnsRecoid(mata.Then.Name, mata.Then.ZoneID)
					if ok {
						log.Printf("获取解析【%s】成功\n", mata.Then.Name)
					} else {
						log.Printf("获取解析【%s】失败\n", mata.Then.Name)
						return
					}
					if dns.Content != mata.Then.Content {
						send = true
						log.Printf("修改解析【%s】\n", mata.Then.Name)
						utils.Dns(mata.Then, dns.ID, mata.Then.ZoneID)
					} else {
						log.Printf("无需修改解析【%s】\n", mata.Main.Name)
					}
				}

				if send {
					log.Println("开始发送通知")
					if conf.Config.BotToken != "" && conf.Config.ChatID != "" {
						msg = "【" + mata.PS + "】" + msg
						go utils.SendTG("#MATA " + msg)
					}
					if conf.Config.ServerJiang != "" {
						go utils.SendSJ(msg)
					}
				} else {
					log.Println("无需发送通知")
				}
			}(mata)
		}

		wg.Wait() // 等待所有 goroutine 完成

		if Once {
			return
		}
		time.Sleep(time.Duration(conf.Config.Corn) * time.Second)
	}
}

func init() {
	once := flag.Bool("once", false, "Run once")
	port := flag.String("port", "8080", "Web port")
	conf.WebPort = *port
	configPath := flag.String("config", "mata.json", "Config file path")
	flag.Parse()

	if *once {
		Once = true
	}
	configInit(*configPath)
}

func configInit(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(path + "配置文件不存在")
		os.Exit(0)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf.Config)
	if err != nil {
		fmt.Println("配置文件错误\n路径" + path + "\n错误" + err.Error())
		os.Exit(0)
	}
	conf.Config.TgApiUrl = strings.TrimRight(conf.Config.TgApiUrl, "/")
}
