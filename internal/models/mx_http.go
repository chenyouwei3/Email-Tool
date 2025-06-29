package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MXResolverHttp struct {
	BaseMxResolver
	ProxyURL  string            //[是否开启代理]代理proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	DohURL    string            //doh的api
	DohParams map[string]string //请求参数
}

func (m MXResolverHttp) mxRecordsTesting(domain string) ([]*net.MX, error) {
	var data []*net.MX
	var err error
	transport := &http.Transport{}
	if len(m.ProxyURL) != 0 {
		proxyUrl, _ := url.Parse(m.ProxyURL)
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	transport.DialContext = (&net.Dialer{
		Timeout: m.BaseMxResolver.ConnTimeout,
	}).DialContext
	client := &http.Client{
		Transport: transport,
	}
	params := url.Values{}
	params.Set("name", domain)
	for param, value := range m.DohParams {
		params.Set(param, value)
	}
	fullURL := m.DohURL + "?" + params.Encode()
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	//处理结果
	var result dnsResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	for _, ans := range result.Answer {
		if ans.Type == 15 { // MX
			parts := strings.Fields(ans.Data)
			pref, _ := strconv.Atoi(parts[0])
			data = append(data, &net.MX{
				Host: parts[1],
				Pref: uint16(pref),
			})
		}
	}
	return data, nil
}

// 表示 DNS 查询的完整响应结构（使用匿名嵌套结构体）
type dnsResponse struct {
	Status int  `json:"Status"` // DNS 响应状态码：0 表示无错误（NOERROR）
	TC     bool `json:"TC"`     // 是否被截断
	RD     bool `json:"RD"`     // 是否启用递归查询
	RA     bool `json:"RA"`     // 服务器是否支持递归
	AD     bool `json:"AD"`     // 响应数据是否通过 DNSSEC 验证
	CD     bool `json:"CD"`     // 是否禁用 DNSSEC 检查

	// 匿名结构体定义查询项
	Question []struct {
		Name string `json:"name"` // 查询的域名
		Type int    `json:"type"` // 查询类型，例如 1=A, 15=MX
	} `json:"Question"`

	// 使用显式结构体定义回答项
	Answer []struct {
		Name string `json:"name"` // 响应域名
		Type int    `json:"type"` // 记录类型
		TTL  int    `json:"TTL"`  // 生存时间
		Data string `json:"data"` // 记录值
	} `json:"Answer"`
}
