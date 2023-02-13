# encapsulation zap

## PART1. 微服务支付模块service开发

- 支付模块ZAP日志工具封装
- 支付模块proto开发
- 支付模块handler开发
- 支付模块main.go开发

## PART2. ZAP开发

关于zap的开发是在common中的.

拉取依赖:

`go get go.uber.org/zap`

`go get gopkg.in/natefinch/lumberjack.v2`    // 用于拆分日志大小

`common/zap.go`:

```go
package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.SugaredLogger
)

// init 用于初始化日志库
func init() {
	// 日志文件名称
	fileName := "micro.log"

	jackLogger := &lumberjack.Logger{
		// 日志文件名称
		Filename: fileName,
		// 日志文件大小 单位:MB
		// 超过该大小将切割日志
		MaxSize: 512,
		// 日志文件时长 单位: 天
		// 超过该时长将切割日志
		//MaxAge:     0,
		// 日志文件的最大备份个数
		MaxBackups: 0,
		// 日志文件中是否使用本地时间
		LocalTime: true,
		// 是否启用压缩
		Compress: true,
	}

	// 将日志写入文件的writer
	syncWriter := zapcore.AddSync(jackLogger)

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 创建编码器
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(
		// 编码器
		jsonEncoder,
		// 写入器
		syncWriter,
		// 日志等级
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	log := zap.New(
		core,
		zap.AddCaller(),
		// 参数和调用层级有关
		zap.AddCallerSkip(1),
	)

	logger = log.Sugar()
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
```