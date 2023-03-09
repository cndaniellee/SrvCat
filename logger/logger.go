package logger

import (
	"SrvCat/config"
	"github.com/kataras/golog"
	"os"
	"time"
)

func init() {
	// 设置日志文件
	logPath := config.Config.LogPath
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		golog.Fatalf("An error occurred while creating the log directory", err)
	}
	logfile := getLogFile(logPath)
	golog.AddOutput(logfile)
	// 设置日志级别
	if config.Config.Debug {
		golog.SetLevel("debug")
	}
}

func getLogFile(logPath string) *os.File {
	today := time.Now().Format("2006-01-02")
	filename := logPath + "/" + today + ".log"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		golog.Fatalf("An error occurred while creating/loading the log file", err)
	}
	return file
}
