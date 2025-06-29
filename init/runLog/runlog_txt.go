package runLog

import (
	conf "email-tool/init/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	zapLog  *zap.Logger
	Control = control{}
)

func InitRunLog(path string) error {
	// 配置日志格式
	config := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	//每次运行的日志文件
	logFilePath := path + time.Now().Format("2006-01-02 15") + ".log"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	writeSyncer := zapcore.AddSync(file)
	// 创建 core
	var core zapcore.Core
	if conf.Conf.APP.Mode == "debug" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			writeSyncer,
			zapcore.DebugLevel, // 开发环境建议是 Debug
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			writeSyncer,
			zapcore.InfoLevel,
		)
	}
	zapLog = zap.New(core)
	Control = control{Zap: zapLog}
	return nil
}

type control struct {
	Zap *zap.Logger
}

func (c *control) Debug(msg string, err error) {
	formatMsg := fmt.Sprintf("%s : %v", msg, err)
	if err != nil {
		fmt.Println(formatMsg)
		c.Zap.Error(formatMsg)
	} else {
		fmt.Println(msg)
		c.Zap.Info(msg)
	}
}

// 邮箱信息,验证状态,错误信息,邮箱是否正确,域名是否正确
func (c *control) EmailsData(email, kind string, err error, isTrueEmail, isTrueDomain bool) {
	// 格式化消息（对齐邮箱，错误信息居中显示）
	formatMsg := fmt.Sprintf(
		"[邮箱: %-35s] [方式: %-5s] [邮箱正确: %-5v] [域名是否合法: %-5v] [错误: %v]",
		email,
		kind,
		isTrueEmail,
		isTrueDomain,
		err,
	)
	// 控制台打印
	fmt.Println(formatMsg)
	// zap 记录日志
	if err != nil {
		c.Zap.Error(formatMsg)
	} else {
		c.Zap.Info(formatMsg)
	}

}
