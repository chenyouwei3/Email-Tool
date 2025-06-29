package models

import (
	"context"
	"net"
)

// 默认解析器使用tcp_udp协议解析  无法使用代理
type MXResolverDefault struct {
	BaseMxResolver
	DnsServer string //是否指定DNS服务器来解析
}

func (m MXResolverDefault) mxRecordsTesting(domain string) ([]*net.MX, error) {
	var data []*net.MX
	var err error
	if len(m.DnsServer) == 0 {
		data, err = net.LookupMX(domain)
		if err != nil {
			return nil, err
		}
	} else {
		dialer := &net.Dialer{
			Timeout: m.BaseMxResolver.ConnTimeout, //连接时间5s
		}
		resolver := &net.Resolver{
			PreferGo: true, //用 Go 内置的 DNS 实现  不是运行操作系统的dns配置
			//   类型func(ctx, network, address) (net.Conn, error)；
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, m.DnsServer)
			},
		}
		data, err = resolver.LookupMX(context.Background(), domain)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
