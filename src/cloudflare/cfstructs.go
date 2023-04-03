package cloudflare

type CFResp struct {
	DNSEntry   []DNSEntry `json:"result"`
	Success    bool       `json:"success"`
	ResultInfo ResultInfo `json:"result_info"`
}

type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type DNSEntry struct {
	ID        string `json:"id"`
	ZoneID    string `json:"zone_id"`
	ZoneName  string `json:"zone_name"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	Locked    bool   `json:"locked"`
	Comment   string `json:"comment"`
	Created   string `json:"created_on"`
	Modified  string `json:"modified_on"`
}
