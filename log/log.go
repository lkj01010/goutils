// Copyright 2015 Chen Xianren. All rights reserved.

// https://studygolang.com/articles/2644

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	LevelDebug = (iota + 1) * 10
	LevelInfo
	LevelWarning
	LevelError
	LevelPanic
	LevelFatal
)

var (
	levels = map[int]string{
		LevelDebug:   "[D]",
		LevelInfo:    "[I]",
		LevelWarning: "[W]",
		LevelError:   "[E]",
		LevelPanic:   "[P]",
		LevelFatal:   "[F]",
	}
)

/*
说明：
前景色            背景色           颜色
---------------------------------------
30                40              黑色
31                41              红色
32                42              绿色
33                43              黃色
34                44              蓝色
35                45              紫红色
36                46              青蓝色
37                47              白色
显示方式           意义
-------------------------
0                终端默认设置
1                高亮显示
4                使用下划线
5                闪烁
7                反白显示
8                不可见

例子：
\033[1;31;40m    <!--1-高亮显示 31-前景色红色  40-背景色黑色-->
\033[0m          <!--采用终端默认设置，即取消颜色设置-->
*/
var (
	colorPrefix = map[int]string{
		LevelDebug:   "\033[34m",
		LevelInfo:    "\033[32m",
		LevelWarning: "\033[1;33m",
		LevelError:   "\033[1;31m",
		LevelPanic:   "\033[1;31m",
		LevelFatal:   "\033[1;31m",
	}
	colorSuffix = "\033[0m"
)

func SetLevelName(level int, name string) {
	levels[level] = name
}

func LevelName(level int) string {
	name, ok := levels[level]
	if !ok {
		name = "LEVEL" + strconv.Itoa(level)
	}
	return name
}

func NameLevel(name string) int {
	for k, v := range levels {
		if v == name {
			return k
		}
	}
	var level int
	if strings.HasPrefix(name, "LEVEL") {
		level, _ = strconv.Atoi(name[5:])
	}
	return level
}

func levelColorPrefix(level int) string {
	p, ok := colorPrefix[level]
	if !ok {
		log.Println("ERROR: levelColorPrefix not found, level=%v", level)
		return ""
	}
	return p
}

func levelColorSuffix() string {
	return colorSuffix
}

type Logger struct {
	mu     sync.Mutex
	level  int
	logger *log.Logger
}

func New(out io.Writer, prefix string, flag, level int) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, prefix, flag),
	}
}

func (self *Logger) Flags() int {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.logger.Flags()
}

func (self *Logger) SetFlags(flag int) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.logger.SetFlags(flag)
}

func (self *Logger) Prefix() string {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.logger.Prefix()
}

func (self *Logger) SetPrefix(prefix string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.logger.SetPrefix(prefix)
}

func (self *Logger) Level() int {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.level
}

func (self *Logger) SetLevel(level int) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.level = level
}

func (self *Logger) Err(level, calldepth int, err error) error {
	if err != nil {
		self.mu.Lock()
		defer self.mu.Unlock()
		if level >= self.level {
			return self.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), err))
		}
	}
	return nil
}

func (self *Logger) ErrDebug(err error) {
	self.Err(LevelDebug, 3, err)
}

func (self *Logger) ErrInfo(err error) {
	self.Err(LevelInfo, 3, err)
}

func (self *Logger) ErrWarning(err error) {
	self.Err(LevelWarning, 3, err)
}

func (self *Logger) ErrError(err error) {
	self.Err(LevelError, 3, err)
}

func (self *Logger) ErrPanic(err error) {
	if err != nil {
		self.Err(LevelPanic, 3, err)
		panic(err)
	}
}

func (self *Logger) ErrFatal(err error) {
	if err != nil {
		self.Err(LevelFatal, 3, err)
		os.Exit(1)
	}
}

//lkj add:
func formatOutput(level int, v ...interface{}) string {
	return fmt.Sprintf("%s%s: %s%s",
		levelColorPrefix(level),
		LevelName(level),
		fmt.Sprint(v...),
		levelColorSuffix())
}

func (self *Logger) Output(level, calldepth int, v ...interface{}) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if level >= self.level {
		return self.logger.Output(calldepth,
			formatOutput(level, v),
		)
	}
	return nil
}

func (self *Logger) Outputf(level, calldepth int, format string, v ...interface{}) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if level >= self.level {
		//lkj modify:
//		return self.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), fmt.Sprintf(format, v...)))
		//-->
		return self.logger.Output(calldepth,
			formatOutput(level, fmt.Sprintf(format, v...)),
		)
		//]]
	}
	return nil
}

func (self *Logger) Outputln(level, calldepth int, v ...interface{}) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	if level >= self.level {
		s := fmt.Sprintln(v...)
		s = s[:len(s) - 1]
//		return self.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), s))
		//-->
		return self.logger.Output(calldepth,
			formatOutput(level, s),
		)
		//]]
	}
	return nil
}

func (self *Logger) Debug(v ...interface{}) {
	self.Output(LevelDebug, 3, v...)
}

func (self *Logger) Info(v ...interface{}) {
	self.Output(LevelInfo, 3, v...)
}

func (self *Logger) Warning(v ...interface{}) {
	self.Output(LevelWarning, 3, v...)
}

func (self *Logger) Error(v ...interface{}) {
	self.Output(LevelError, 3, v...)
}

func (self *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	self.Output(LevelPanic, 3, s)
	panic(s)
}

func (self *Logger) Fatal(v ...interface{}) {
	self.Output(LevelFatal, 3, v...)
	os.Exit(1)
}

func (self *Logger) Debugf(format string, v ...interface{}) {
	self.Outputf(LevelDebug, 3, format, v...)
}

func (self *Logger) Infof(format string, v ...interface{}) {
	self.Outputf(LevelInfo, 3, format, v...)
}

func (self *Logger) Warningf(format string, v ...interface{}) {
	self.Outputf(LevelWarning, 3, format, v...)
}

func (self *Logger) Errorf(format string, v ...interface{}) {
	self.Outputf(LevelError, 3, format, v...)
}

func (self *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	self.Outputf(LevelPanic, 3, "%s", s)
	panic(s)
}

func (self *Logger) Fatalf(format string, v ...interface{}) {
	self.Outputf(LevelFatal, 3, format, v...)
	os.Exit(1)
}

func (self *Logger) Debugln(v ...interface{}) {
	self.Outputln(LevelDebug, 3, v...)
}

func (self *Logger) Infoln(v ...interface{}) {
	self.Outputln(LevelInfo, 3, v...)
}

func (self *Logger) Warningln(v ...interface{}) {
	self.Outputln(LevelWarning, 3, v...)
}

func (self *Logger) Errorln(v ...interface{}) {
	self.Outputln(LevelError, 3, v...)
}

func (self *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	s = s[:len(s) - 1]
	self.Outputln(LevelPanic, 3, s)
	panic(s)
}

func (self *Logger) Fatalln(v ...interface{}) {
	self.Outputln(LevelFatal, 3, v...)
	os.Exit(1)
}

//lkj: set level to DEBUG
var std = New(os.Stderr, "", log.LstdFlags | log.Lshortfile, LevelDebug)

func SetOutput(w io.Writer) {
	*std = *New(w, std.logger.Prefix(), std.logger.Flags(), std.level)
}

func Flags() int {
	return std.Flags()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func Prefix() string {
	return std.Prefix()
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func Level() int {
	return std.Level()
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func ErrDebug(err error) {
	std.Err(LevelDebug, 3, err)
}

func ErrInfo(err error) {
	std.Err(LevelInfo, 3, err)
}

func ErrWarning(err error) {
	std.Err(LevelWarning, 3, err)
}

func ErrError(err error) {
	std.Err(LevelError, 3, err)
}

func ErrPanic(err error) {
	if err != nil {
		std.Err(LevelPanic, 3, err)
		panic(err)
	}
}

func ErrFatal(err error) {
	if err != nil {
		std.Err(LevelFatal, 3, err)
		os.Exit(1)
	}
}

// note: 这几个函数输出会被[]包住，原因是Sprint里判断是非string，是interface{}
func Debug(v ...interface{}) {
	std.Output(LevelDebug, 3, v...)
}

func Info(v ...interface{}) {
	std.Output(LevelInfo, 3, v...)
}

func Warning(v ...interface{}) {
	std.Output(LevelWarning, 3, v...)
}

func Error(v ...interface{}) {
	std.Output(LevelError, 3, v...)
}
//]]

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(LevelPanic, 3, s)
	panic(s)
}

func Fatal(v ...interface{}) {
	std.Output(LevelFatal, 3, v...)
	os.Exit(1)
}

func Debugf(format string, v ...interface{}) {
	std.Outputf(LevelDebug, 3, format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Outputf(LevelInfo, 3, format, v...)
}

func Warningf(format string, v ...interface{}) {
	std.Outputf(LevelWarning, 3, format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Outputf(LevelError, 3, format, v...)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Outputf(LevelPanic, 3, "%s", s)
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	std.Outputf(LevelFatal, 3, format, v...)
	os.Exit(1)
}

func Debugln(v ...interface{}) {
	std.Outputln(LevelDebug, 3, v...)
}

func Infoln(v ...interface{}) {
	std.Outputln(LevelInfo, 3, v...)
}

func Warningln(v ...interface{}) {
	std.Outputln(LevelWarning, 3, v...)
}

func Errorln(v ...interface{}) {
	std.Outputln(LevelError, 3, v...)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	s = s[:len(s) - 1]
	std.Outputln(LevelPanic, 3, s)
	panic(s)
}

func Fatalln(v ...interface{}) {
	std.Outputln(LevelFatal, 3, v...)
	os.Exit(1)
}

