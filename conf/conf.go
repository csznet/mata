package conf

// DNSRecord 用于构造 DNS 记录的 JSON 结构体
type DNSRecord struct {
	ZoneID  string `json:"ZoneID"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

type RecoidRes struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result []OneRes `json:"result"`
}

type OneRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Mata struct {
	Target string
	PS     string `json:"ps"`
	Main   DNSRecord
	Then   DNSRecord
}

type cf struct {
	ApiKey   string `json:"cf_api_key"`
	ZoneID   string
	BotToken string
	ChatID   string
	TgApiUrl string
	Corn     int64
	Mata     []Mata
}

var Config cf
