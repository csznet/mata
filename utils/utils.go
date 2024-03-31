package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"csz.net/mata/conf"
)

func Check(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false // 连接失败
	}
	defer conn.Close() // 确保关闭连接
	return true        // 连接成功
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

func SendMessage(text string) {
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
