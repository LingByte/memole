package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// 颜色代码
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// Logger 日志器结构
type Logger struct {
	level LogLevel
}

// NewLogger 创建新的日志器
func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel 获取日志级别
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// getLevelString 获取级别字符串
func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// getLevelColor 获取级别颜色
func (l *Logger) getLevelColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return Cyan
	case INFO:
		return Green
	case WARN:
		return Yellow
	case ERROR:
		return Red
	case FATAL:
		return Red + Bold
	default:
		return White
	}
}

// log 内部日志方法
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := l.getLevelString(level)
	levelColor := l.getLevelColor(level)
	
	message := fmt.Sprintf(format, args...)
	
	// 带颜色的日志输出
	logLine := fmt.Sprintf("%s[%s]%s %s %s%s%s\n", 
		levelColor, levelStr, Reset, 
		timestamp, 
		levelColor, message, Reset)
	
	fmt.Print(logLine)
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// 全局日志器实例
var GlobalLogger = NewLogger(INFO)

// 全局函数，方便直接调用
func Debug(format string, args ...interface{}) {
	GlobalLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GlobalLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GlobalLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GlobalLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GlobalLogger.Fatal(format, args...)
}

func SetLevel(level LogLevel) {
	GlobalLogger.SetLevel(level)
}

func GetLevel() LogLevel {
	return GlobalLogger.GetLevel()
}
