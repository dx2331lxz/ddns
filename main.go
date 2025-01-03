package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const APITOKEN = "XXXX"

type Record struct {
	LoginTime  string `json:"login_time"`
	IsOwnerIP  string `json:"is_owner_ip"`
	LoginIP    string `json:"login_ip"`
	LoginIPV6  string `json:"login_ip_v6"`
	MacAddress string `json:"mac_address"`
	PhoneFlag  int    `json:"phone_flag"`
	PlayButton string `json:"play_button"`
}

type Response struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Count   int      `json:"count"`
	Records []Record `json:"records"`
}

//	curl 'https://xha.ouc.edu.cn:802/eportal/portal/page/loadOnlineRecord?callback=dr1004&lang=zh-CN&program_index=ctshNw1713845951&page_index=V5fmKw1713845966&user_account=22010022045&wlan_user_ip=0.0.0.0&wlan_user_mac=000000000000&start_time=2010-01-01&end_time=2100-01-01&start_rn=1&end_rn=5&jsVersion=4.1&v=4008&lang=zh' \
//	 -H 'accept: */*' \
//	 -H 'accept-language: zh-CN,zh;q=0.9' \
//	 -H 'cache-control: no-cache' \
//	 -H 'cookie: identifyId=d90f9a462eef4cc1bf865aatywerwe23' \
//	 -H 'pragma: no-cache' \
//	 -H 'referer: https://xha.ouc.edu.cn/' \
//	 -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
//	 -H 'sec-ch-ua-mobile: ?0' \
//	 -H 'sec-ch-ua-platform: "macOS"' \
//	 -H 'sec-fetch-dest: script' \
//	 -H 'sec-fetch-mode: no-cors' \
//	 -H 'sec-fetch-site: same-site' \
//	 -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'
func fetchOnlineRecord(user_account string) string {
	// 基础 URL
	baseURL := "https://xha.ouc.edu.cn:802/eportal/portal/page/loadOnlineRecord"

	// URL 参数
	params := url.Values{}
	params.Add("callback", "dr1004")
	params.Add("lang", "zh-CN")
	params.Add("program_index", "ctshNw1713845951")
	params.Add("page_index", "V5fmKw1713845966")
	params.Add("user_account", user_account)
	params.Add("wlan_user_ip", "0.0.0.0")
	params.Add("wlan_user_mac", "000000000000")
	params.Add("start_time", "2010-01-01")
	params.Add("end_time", "2100-01-01")
	params.Add("start_rn", "1")
	params.Add("end_rn", "5")
	params.Add("jsVersion", "4.1")
	params.Add("v", "4008")

	// 动态构建完整 URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}

	// 设置请求头
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "identifyId=d90f9a462eef4cc1bf865aatywerwe23")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://xha.ouc.edu.cn/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "script")
	req.Header.Set("sec-fetch-mode", "no-cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}

	// 打印响应

	return string(body)
}

// 解析出login_ip
func parseLoginIP(body string) string {
	// 提取 JSON 内容
	jsonData := strings.TrimSuffix(strings.TrimPrefix(body, "dr1004("), ");")

	// 解析 JSON 数据
	var response Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("JSON 解析失败:", err)
		return ""
	}

	// 查找 is_owner_ip 为 1 的记录
	for _, record := range response.Records {
		if record.IsOwnerIP == "1" {
			fmt.Println("找到的 login_ip:", record.LoginIP)
			return record.LoginIP
		}
	}

	fmt.Println("未找到 is_owner_ip 为 1 的记录")
	return ""
}

// Zone represents the zone information structure
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ZoneResponse represents the Cloudflare API response structure
type ZoneResponse struct {
	Result []Zone `json:"result"`
}

// 获取zoon_id
func getZoneID(domain string) string {
	baseURL := "https://api.cloudflare.com/client/v4/zones"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+APITOKEN)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {

	}
	// 打印响应
	var apiResponse ZoneResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return ""
	}
	for _, zone := range apiResponse.Result {
		if zone.Name == domain {
			return zone.ID
		}
	}
	fmt.Println("未找到 zone_id")
	return ""
}

type DNSRecord struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

// 获取DNS记录
func getDNSRecord(zoneID string, DnsRecord string) string {
	baseURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+APITOKEN)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}
	//fmt.Println(string(body))
	var apiResponse DNSRecord
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return ""
	}
	for _, record := range apiResponse.Result {
		if record.Name == DnsRecord {
			return record.ID
		}
	}
	fmt.Println("未找到 DNS 记录")
	return ""
}

// Overwrite DNS Record
func overwriteDNSRecord(zoneID string, recordID string, name string, record string) {
	baseURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	data := map[string]interface{}{
		"type":    "A",
		"name":    name,
		"content": record,
		"ttl":     60,
		"proxied": false,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}
	req, err := http.NewRequest("PUT", baseURL, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+APITOKEN)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	//	提取 JSON 内容中的success字段
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON 解析失败:", err)
		return
	}
	if response["success"] == true {
		fmt.Println("DNS 记录更新成功")
	} else {
		fmt.Println("DNS 记录更新失败")
	}
}

func main() {
	// 定义学号
	userAccount := "XXX"
	// 域名
	domain := "XXX.com"
	// 记录
	subDomain := "XXX.XXX.com"

	// 获取在线记录
	body := fetchOnlineRecord(userAccount)
	if body == "" {
		fmt.Println("获取在线记录失败")
		return
	}
	loginIp := parseLoginIP(body)
	if loginIp == "" {
		fmt.Println("解析 login_ip 失败")
		return
	}
	fmt.Println("login_ip:", loginIp)
	ZoneID := getZoneID(domain)
	if ZoneID == "" {
		fmt.Println("获取 zone_id 失败")
		return
	}
	fmt.Println("zone_id:", ZoneID)
	record := getDNSRecord(ZoneID, subDomain)
	if record == "" {
		fmt.Println("获取 DNS 记录失败")
		return
	}
	fmt.Println("DNS 记录 ID:", record)
	overwriteDNSRecord(ZoneID, record, subDomain, loginIp)
}
