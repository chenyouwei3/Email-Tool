package emailData

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 初始化邮箱数据
func initEmailData(path string) ([]string, error) {
	var data []string
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//使用scanner读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "@", 2)
		SystemCensus[parts[1]]++
		data = append(data, line)
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// 初始化黑白名单
func initEmailWhiteBlackList(path string, isTrue bool) error {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	//使用scanner读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		SystemCensus[line]++
		SystemBlackWhiteList[line] = isTrue

	}
	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

// 输出数组到txt文件
func EmailsOutput(filename string, emails []string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer file.Close()

	// 使用缓冲写入器提高性能
	writer := bufio.NewWriter(file)
	for _, line := range emails {
		_, err := writer.WriteString(line + "\n") // 每个元素一行
		if err != nil {
			fmt.Println("写入失败:", err)
			return
		}
	}
	writer.Flush() // 写入缓冲区内容到文件
}
