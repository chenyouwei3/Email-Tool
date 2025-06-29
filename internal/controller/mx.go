package controller

import (
	conf "email-tool/init/config"
	"email-tool/internal/models"
	"sync"
)

func mxController(domains map[string]int, blackWhiteList models.BlackWhiteList, systemBlackWhiteList map[string]bool, workers int) (retErr error) {
	var mxResolver models.MxResolver
	baseMxResolver := models.BaseMxResolver{
		ConnTimeout: conf.Conf.APP.MxConfig.Timeout,
	}
	switch conf.Conf.APP.MxConfig.Mode {
	case "doh":
		mxResolver = &models.MXResolverHttp{ //使用google
			BaseMxResolver: baseMxResolver,
			ProxyURL:       "http://127.0.0.1:7890",
			DohURL:         "https://dns.google/resolve",
			DohParams: map[string]string{
				"type": "MX",
			},
		}
	default:
		mxResolver = &models.MXResolverDefault{
			BaseMxResolver: baseMxResolver,
			DnsServer:      conf.Conf.APP.MxConfig.DnsServer,
		}
	}
	//----------并发准备-----------------
	var (
		semaphore = make(chan struct{}, workers)
		wg        sync.WaitGroup
		errOnce   sync.Once
	)
	//------------------------------
	for domain := range domains {
		//黑名单白名单都没有必要进行mx验证
		if _, ok := systemBlackWhiteList[domain]; ok {
			continue
		}
		wg.Add(1)
		semaphore <- struct{}{} //限制协程量
		go func(domain string) {
			defer wg.Done()
			defer func() {
				<-semaphore // 释放信号量位置
			}()
			mxRecords, err := models.MxBegin(mxResolver, domain)
			if err != nil {
				errOnce.Do(func() {
					retErr = err
				})
				mxRecords, err = models.MxBegin(mxResolver, domain)
			}

			blackWhiteList.Mutex.Lock()
			if len(mxRecords) != 0 {
				mxRecordsChan <- mxRecordsDefault{
					Domain: domain,
					Data:   mxRecords,
				}
				blackWhiteList.Data[domain] = true
			} else {
				blackWhiteList.Data[domain] = false
			}
			blackWhiteList.Mutex.Unlock()
		}(domain)
	}
	wg.Wait()
	return nil
}
