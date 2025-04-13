package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"csz.net/mata/conf"
)

// imcp check
func Ping(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("ip4:icmp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func Tcp(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false // 连接失败
	}
	defer conn.Close() // 确保关闭连接
	return true        // 连接成功
}
func Check(address string, timeout time.Duration) (bool, bool) {
	check, status := false, false
	// 判断是否由远程进行检测
	if len(conf.Config.Contcp) > 5 {
		// 检查是否有/结尾
		if conf.Config.Contcp[len(conf.Config.Contcp)-1:] != "/" {
			conf.Config.Contcp += "/"
		}
		// 发送请求
		resp, err := http.Get(conf.Config.Contcp + address)
		if err != nil {
			return false, false
		}
		defer resp.Body.Close()
		// 获取body内容
		var body bytes.Buffer
		_, err = io.Copy(&body, resp.Body)
		if err != nil {
			return false, false
		}
		check = true
		if body.String() == "true" {
			status = true
		}
		return check, status
	}
	check = true
	if strings.ContainsAny(address, ":") {
		return check, Tcp(address, timeout)
	}
	return check, Ping(address, timeout)
}

func Dns(record conf.DNSRecord, recordID string, ZoneID string) bool {
	url := "https://api.cloudflare.com/client/v4/zones/" + ZoneID + "/dns_records/" + recordID
	recordBytes, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Error marshalling DNS record:", err)
		return false
	}
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(recordBytes))
	if err != nil {
		return false
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+conf.Config.ApiKey)
	http.DefaultClient.Do(req)
	return true
}

func GetDnsRecoid(recoid string, ZoneID string) (bool, conf.OneRes) {
	var no conf.OneRes
	url := "https://api.cloudflare.com/client/v4/zones/" + ZoneID + "/dns_records"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, no
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+conf.Config.ApiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, no
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, no
	}
	var result conf.RecoidRes
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, no
	}
	if !result.Success {
		return false, no
	}
	for _, record := range result.Result {
		if record.Name == recoid {
			return true, record
		}
	}
	return false, no
}

func SendTG(text string) {
	url := fmt.Sprintf(conf.Config.TgApiUrl+"/bot%s/sendMessage", conf.Config.BotToken)
	requestBody := fmt.Sprintf(`{"chat_id": "%s", "text": "%s"}`, conf.Config.ChatID, text)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: ", resp.Status)
	} else {
		fmt.Println("Message sent successfully")
	}
}

// SendSJ 发送到Server酱
func SendSJ(text string) {
	resp, err := http.Get("https://sc.ftqq.com/" + conf.Config.ServerJiang + ".send?text=" + text)
	if err != nil {
		fmt.Println("Error sending request: ", err)
	}
	defer resp.Body.Close()
}

func Web() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		node := r.URL.Query().Get("node")
		if node != "" {
			// 更新Array的值为当前时间
			conf.Array[node] = int(time.Now().Unix())
			log.Println("更新节点", node, "最后存活时间", time.Unix(int64(conf.Array[node]), 0).Format("2006-01-02 15:04:05"))
			fmt.Fprintf(w, "ok")
		}
	})
	http.ListenAndServe(":"+conf.WebPort, nil)
}
