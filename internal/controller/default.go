package controller

import (
	conf "email-tool/init/config"
	"email-tool/init/emailData"
	"email-tool/init/runLog"
	"email-tool/internal/models"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	trueEmailChan  = make(chan string, 25) //验证确定数据
	trueEmailData  []string
	falseEmailChan = make(chan string, 25) //验证失败数据
	falseEmailData []string
	mxRecordsChan  = make(chan mxRecordsDefault, 25)
)

type mxRecordsDefault struct {
	Data   []*net.MX
	Domain string
}

func ReceiveData() {
	for {
		select {
		case email := <-trueEmailChan:
			trueEmailData = append(trueEmailData, email)
		case email := <-falseEmailChan:
			falseEmailData = append(falseEmailData, email)
		case mxRecord := <-mxRecordsChan:
			emailData.DomainMxRecords.Mutex.Lock()
			emailData.DomainMxRecords.Data[mxRecord.Domain] = mxRecord.Data
			emailData.DomainMxRecords.Mutex.Unlock()
		}
	}
}

func DefaultRouter(regex, mx, smtp bool, emails []string) {
	go ReceiveData() //协程接收数据
	if regex {
		testing := regexController(emails, conf.Conf.APP.RegexConfig.Goroutines)
		fmt.Println("正则表达式通过率", len(testing)/len(emails))
	}
	if mx {
		if err := mxController(
			emailData.SystemCensus,
			emailData.DomainBlackWhiteList,
			emailData.SystemBlackWhiteList,
			10,
		); err != nil {
			runLog.Control.Debug("mx记录检测失败", err)
			return
		}
	}

	if smtp {
		fmt.Println("===========开始smtp验证============")
		var smtpResolver models.SmtpResolver
		smtpResolver = &models.SmtpResolverDefault{
			BaseSmtpResolver: models.BaseSmtpResolver{
				ConnTimeout: conf.Conf.APP.SmtpConfig.Timeout,
				MxRecords:   emailData.DomainMxRecords.Data,
				Port:        conf.Conf.APP.SmtpConfig.Port,
			},
		}
		var (
			semaphore = make(chan struct{}, conf.Conf.APP.SmtpConfig.Goroutines)
			wg        sync.WaitGroup
		)
		for _, email := range emails {
			parts := strings.SplitN(email, "@", 2)
			domain := parts[1]
			//对黑白名单的进行特殊处理
			if isTrue, ok := emailData.SystemBlackWhiteList[domain]; ok {
				if isTrue {
					runLog.Control.EmailsData(email, "白名单内的数据", nil, true, true)
					trueEmailChan <- email
				} else {
					runLog.Control.EmailsData(email, "黑名单内的数据", nil, false, false)
					falseEmailChan <- email
				}
				continue
			}
			// 使用解析出来的黑名单
			if value, ok := emailData.DomainBlackWhiteList.Data[domain]; ok && !value {
				runLog.Control.EmailsData(email, "解析出来域名错误", nil, false, false)
				falseEmailChan <- email
				continue
			}
			wg.Add(1)
			semaphore <- struct{}{} // 阻塞直到有空位
			go func(email, domain string) {
				defer wg.Done()
				defer func() {
					<-semaphore
				}()
				if err := models.SmtpBegin(smtpResolver, email); err != nil {
					errStr := strings.ToLower(err.Error())
					isMatched := false //是否匹配
					for _, errMsg := range conf.Conf.APP.SmtpConfig.Keys {
						if strings.Contains(errStr, errMsg) {
							runLog.Control.EmailsData(email, "关键词之内", err, false, false)
							falseEmailChan <- email
							isMatched = true
							break
						}
					}
					if !isMatched {
						runLog.Control.EmailsData(email, "关键词范围之外", err, true, true)
						trueEmailChan <- email
					}
				} else {
					runLog.Control.EmailsData(email, "验证存在", nil, true, true)
					trueEmailChan <- email
				}
			}(email, domain)
		}
		wg.Wait()
	}
	fmt.Println("===========检测完毕,正在写入数据==============")
	emailData.EmailsOutput("data/true.txt", trueEmailData)
	emailData.EmailsOutput("data/false.txt", falseEmailData)
	fmt.Println("检测完毕")
	trueEmailData, falseEmailData = nil, nil
}
