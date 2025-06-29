package models

import (
	"net"
	"time"
)

// 默认interface 规定所有的mx测试方法的模板
type MxResolver interface {
	mxRecordsTesting(string) ([]*net.MX, error)
}

type BaseMxResolver struct {
	ConnTimeout time.Duration //发起请求超时
}

func MxBegin(resolver MxResolver, domain string) ([]*net.MX, error) {
	mxRecords, err := resolver.mxRecordsTesting(domain)
	if err != nil {
		return nil, err
	}
	return mxRecords, nil
}
