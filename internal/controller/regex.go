package controller

import (
	conf "email-tool/init/config"
	"regexp"
	"sync"
)

func regexController(emails []string, workers int) []string {
	var resEmails []string

	var mu sync.Mutex

	emailChan := make(chan string, 1000)
	wg := sync.WaitGroup{}
	re := regexp.MustCompile(conf.Conf.APP.RegexConfig.Regex)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for email := range emailChan {
				//验证邮箱
				if re.MatchString(email) {
					mu.Lock()
					resEmails = append(resEmails, email)
					mu.Unlock()
				}
			}
		}()
	}

	for _, email := range emails {
		emailChan <- email
	}
	close(emailChan)
	wg.Wait()

	return resEmails
}
