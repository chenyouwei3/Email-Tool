package models

import (
	"fmt"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SmtpResolverDefault struct {
	BaseSmtpResolver
}

func (s SmtpResolverDefault) smtpConnTesting(email string) error {
	// 使用第一个 MX 服务器
	parts := strings.SplitN(email, "@", 2)

	mxHost := s.MxRecords[parts[1]][0].Host

	addr := fmt.Sprintf("%s:%s", mxHost, s.Port)

	// 建立 TCP 连接
	conn, err := net.DialTimeout("tcp", addr, s.ConnTimeout)
	if err != nil {
		return fmt.Errorf("failed to connect to mail server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return fmt.Errorf("failed to create smtp client: %v", err)
	}
	defer client.Quit()
	sender, domain := generateJapaneseStyleEmail()
	// 使用一个假的发件人地址（通常是你的域名中的地址）
	from := sender
	to := email

	// 发起 SMTP 会话
	if err = client.Hello(domain); err != nil {
		return fmt.Errorf("HELO failed: %v", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %v", err) // 说明目标邮箱可能不存在
	}
	return nil
}

var (
	japaneseSurnames = []string{"tanaka", "suzuki", "yamada", "kobayashi", "saito", "kato", "ito", "fujita", "shimizu", "nakamura"}
	japaneseNames    = []string{"haruki", "yuki", "kenta", "naoki", "hiroshi", "ayumi", "yui", "rin", "takumi", "mei"}
	domainSuffixes   = []string{"docomo.ne.jp", "ezweb.ne.jp", "yahoo.co.jp", "gmail.com", "softbank.ne.jp"}
)

func generateJapaneseStyleEmail() (sender string, domain string) {
	rand.Seed(time.Now().UnixNano())
	surname := japaneseSurnames[rand.Intn(len(japaneseSurnames))]
	name := japaneseNames[rand.Intn(len(japaneseNames))]
	domain = domainSuffixes[rand.Intn(len(domainSuffixes))]
	// 可选连接符："."、"_" 或空
	sep := []string{"", ".", "_"}[rand.Intn(3)]

	// 加一些数字，比如年份、生日等
	num := fmt.Sprintf("%02d", rand.Intn(100))      // 00 ~ 99
	num2 := fmt.Sprintf("%04d", 1980+rand.Intn(30)) // 1980 ~ 2009

	formats := []string{
		fmt.Sprintf("%s%s%s", surname, sep, name),
		fmt.Sprintf("%s%s%s%s", surname, sep, name, num),
		fmt.Sprintf("%s%s%s%s", surname, sep, name, num2),
		fmt.Sprintf("%c.%s%s", surname[0], name, num),
	}

	local := formats[rand.Intn(len(formats))]
	sender = fmt.Sprintf("%s@%s", local, domain)
	return sender, domain
}
