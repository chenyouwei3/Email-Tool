package models

import (
	"net"
	"time"
)

type SmtpResolver interface {
	smtpConnTesting(string) error
}

type BaseSmtpResolver struct {
	ConnTimeout time.Duration        //连接超时
	MxRecords   map[string][]*net.MX //mx记录
	Sender      string               //发件人
	Domain      string               //发件方域名
	Port        string               //smtp 服务器端口
}

func SmtpBegin(resolver SmtpResolver, email string) error {
	err := resolver.smtpConnTesting(email)
	if err != nil {
		return err
	}
	return nil
}
