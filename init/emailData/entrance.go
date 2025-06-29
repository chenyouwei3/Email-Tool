package emailData

import (
	conf "email-tool/init/config"
)

func EmailBegin() (emails []string, err error) {
	//var (
	//	errChan              = make(chan error, 3) // 缓冲防阻塞
	//	once                 sync.Once
	//	wg                   sync.WaitGroup
	//	whiteList, blackList []string
	//)
	//wg.Add(3)
	////读取主邮件数据
	//go func() {
	//	defer wg.Done()
	//	var e error
	//	emails, e = initEmailData(config.Conf.Data.Path)
	//	if e != nil {
	//		errChan <- e
	//		return
	//	}
	//}()
	//// 并发处理白名单
	//go func() {
	//	defer wg.Done()
	//	var e error
	//	whiteList, e = initEmailWhiteBlackList(config.Conf.Data.WhiteList, true)
	//	if e != nil {
	//		errChan <- e
	//		return
	//	}
	//}()
	//
	//// 并发处理黑名单
	//go func() {
	//	defer wg.Done()
	//	var e error
	//	blackList, e = initEmailWhiteBlackList(config.Conf.Data.WhiteList, false)
	//	if e != nil {
	//		errChan <- e
	//		return
	//	}
	//}()
	//
	//wg.Wait()
	//close(errChan)
	//// 检查并发过程中是否有错误
	//for e := range errChan {
	//	once.Do(func() {
	//		err = e
	//	})
	//}
	//读取主邮件数据
	emails, err = initEmailData(conf.Conf.Data.Path)
	if err != nil {
		return
	}
	// 处理黑名单
	err = initEmailWhiteBlackList(conf.Conf.Data.BlackList, false)
	if err != nil {
		return
	}
	// 处理白名单
	err = initEmailWhiteBlackList(conf.Conf.Data.WhiteList, true)
	if err != nil {
		return
	}
	return
}
