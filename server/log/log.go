// Package log 日志打印包，打印到日志文件
package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type myLog struct {
	// 日志类
	logger *log.Logger
	// 日志文件名称
	logFile string
	// 日志文件
	fileWriter *os.File
}

var mylogger myLog

var lock sync.Mutex

func init() {
	var err error
	mylogger.logFile = fmt.Sprintf("./app_data/log/%v.log", time.Now().Format("20060102"))
	mylogger.fileWriter, err = os.OpenFile(mylogger.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	mylogger.logger = log.New(mylogger.fileWriter, "", log.LstdFlags)
	// 新建定时任务，默认每天0点生成新的日志文件
	c := cron.New()
	_, err = c.AddFunc("5 0 * * *", ChangeLogFile)
	if err != nil {
		Errorf("启动定时器失败! ERROR:%v", err)
	}
	c.Start()
}

// changeLogFile 每日凌晨切换到新的日志文件
func ChangeLogFile() {
	lock.Lock()
	defer lock.Unlock()
	newFile := fmt.Sprintf("./log/%v.log", time.Now().Format("20060102"))
	newFileWriter, err := os.OpenFile(mylogger.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Errorf("新建日志文件%s.log错误! ERROR:%v", newFile, err)
		return
	}
	// 切换Logger指向新的日志
	mylogger.logger = log.New(newFileWriter, "", log.LstdFlags)
	// 关闭上次的日志文件
	if mylogger.fileWriter.Close(); err != nil {
		Errorf("关闭日志文件%s.log错误! ERROR:%v", mylogger.logFile, err)
	}
	mylogger.logFile = newFile
	mylogger.fileWriter = newFileWriter

}

// rebuildIfNotExists 如果日志文件不存在就生成新的
func rebuildIfNotExists() {
	_, err := os.Stat(mylogger.logFile)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		ChangeLogFile()
	}
	Errorf("log file exists but stat err:%v", err)
}

// Errorf 错误日志
func Errorf(format string, args ...interface{}) {
	rebuildIfNotExists()
	_, file, line, _ := runtime.Caller(1)
	mylogger.logger.Printf(fmt.Sprintf("%s:%d ERR: %s", file, line, format), args...)
}

// Infof info日志
func Infof(format string, args ...interface{}) {
	rebuildIfNotExists()
	_, file, line, _ := runtime.Caller(1)
	mylogger.logger.Printf(fmt.Sprintf("%s:%d INFO: %s", file, line, format), args...)
}

// Debugf 调试日志
func Debugf(format string, args ...interface{}) {
	rebuildIfNotExists()
	_, file, line, _ := runtime.Caller(1)
	mylogger.logger.Printf(fmt.Sprintf("%s:%d DEBUG: %s", file, line, format), args...)
}
