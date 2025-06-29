package emailData

import (
	"email-tool/internal/models"
	"net"
	"sync"
)

var (
	SystemCensus = make(map[string]int) //统计域名出现了多少次
	//DomainBlackWhiteList = make(map[string]bool) //黑白名单
	//DomainMxRecords      = make(map[string][]*net.MX)
	SystemBlackWhiteList = make(map[string]bool) //黑白名单
)

var (
	DomainBlackWhiteList = models.BlackWhiteList{
		Mutex: &sync.RWMutex{},
		Data:  make(map[string]bool),
	}

	DomainMxRecords = models.MxRecords{
		Mutex: &sync.RWMutex{},
		Data:  make(map[string][]*net.MX),
	}
)
