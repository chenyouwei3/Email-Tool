package main

import (
	"email-tool/init/config"
	"email-tool/init/emailData"
	"email-tool/init/runLog"
	"email-tool/internal/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strconv"
	"time"
)

func init() {
	//初始化配置文件
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	//设置运行模式
	if config.Conf.APP.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	//设置运行日志
	if err = runLog.InitRunLog("../logs/"); err != nil {
		panic(err)
	}
}

func main() {
	runLog.Control.Debug("配置项加载完毕", nil)
	fmt.Println("  _____                 _ _     _____           _ \n | ____|_ __ ___   __ _(_) |   |_   _|__   ___ | |\n |  _| | '_ ` _ \\ / _` | | |_____| |/ _ \\ / _ \\| |\n | |___| | | | | | (_| | | |_____| | (_) | (_) | |\n |_____|_| |_| |_|\\__,_|_|_|     |_|\\___/ \\___/|_|\n                                                  ")
	defer runLog.Control.Zap.Sync() //运行日志退出
	//数据读取[邮箱数据]-----[域名数据]-----[域名黑白名单]
	emails, err := emailData.EmailBegin()
	if err != nil {
		runLog.Control.Debug("邮箱数据读取失败", err)
		return
	}
	runLog.Control.Debug("邮箱数据读取完毕,邮箱数量:"+strconv.Itoa(len(emails)), nil)
	//数据排序分片并发处理
	/* ---------------------模式选择----------------*/
	// 切片处理，每片 100 个邮箱
	chunks := chunkEmails(emails, 10)
	//系统已默认配置宽带账号
	//请以adsl-start和adsl-stop
	//或pppoe-start和pppoe-stop
	//或/sbin/ifup ppp0和/sbin/ifdown ppp0
	for index, chunk := range chunks {
		fmt.Println("====正在运行第一个数据分片===========", index+1)
		time.Sleep(time.Second * 5)
		err = dialTool("ip", []string{"address"})
		if err != nil {
			runLog.Control.Debug("ip address 指令", nil)
		}
		time.Sleep(time.Second * 5)
		err = dialTool("pppoe-start", nil)
		if err != nil {
			runLog.Control.Debug("pppoe-start 指令", nil)
		}
		time.Sleep(time.Second * 10)
		switch config.Conf.APP.VerifyMode {
		case "smtp":
			controller.DefaultRouter(true, true, true, chunk)
		case "mx":
			controller.DefaultRouter(true, true, false, chunk)
		case "regex":
			controller.DefaultRouter(true, false, false, chunk)
		default:
			controller.DefaultRouter(true, true, true, chunk)
		}
		err = dialTool("pppoe-stop", nil)
		if err != nil {
			runLog.Control.Debug("pppoe-stop 指令", nil)
		}
		time.Sleep(time.Second * 10)
	}
}
func dialTool(cmmd string, args []string) error {
	maxRetry := 5
	for i := 0; i < maxRetry; i++ {
		cmd := exec.Command(cmmd, args...)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("命令 [%s %v] 执行失败(%d/%d): %v\n", cmmd, args, i+1, maxRetry, err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("命令 [%s %v] 执行成功:\n%s\n", cmmd, args, string(output))
		return nil
	}
	return fmt.Errorf("命令 [%s %v] 重试 %d 次仍失败", cmmd, args, maxRetry)
}

// emails是被分片   chunkSize 分片大小
func chunkEmails(emails []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(emails); i += chunkSize {
		end := i + chunkSize
		if end > len(emails) {
			end = len(emails)
		}
		chunk := emails[i:end]
		chunks = append(chunks, chunk)

	}
	return chunks
}
