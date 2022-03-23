// Copyright 2015 Chen Xianren. All rights reserved.

// https://studygolang.com/articles/2644

package log

import (
    "fmt"
    "io"
    "log"
    "os"
    "runtime"
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

func (l *Logger) Flags() int {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.logger.Flags()
}

func (l *Logger) SetFlags(flag int) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.SetFlags(flag)
}

func (l *Logger) Prefix() string {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.logger.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.SetPrefix(prefix)
}

func (l *Logger) Level() int {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.level
}

func (l *Logger) SetLevel(level int) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.level = level
}

func (l *Logger) Err(level, calldepth int, err error) error {
    if err != nil {
        l.mu.Lock()
        defer l.mu.Unlock()
        if level >= l.level {
            return l.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), err))
        }
    }
    return nil
}

func (l *Logger) ErrDebug(err error) {
    l.Err(LevelDebug, 3, err)
}

func (l *Logger) ErrInfo(err error) {
    l.Err(LevelInfo, 3, err)
}

func (l *Logger) ErrWarning(err error) {
    l.Err(LevelWarning, 3, err)
}

func (l *Logger) ErrError(err error) {
    l.Err(LevelError, 3, err)
}

func (l *Logger) ErrPanic(err error) {
    if err != nil {
        l.Err(LevelPanic, 3, err)
        panic(err)
    }
}

func (l *Logger) ErrFatal(err error) {
    if err != nil {
        l.Err(LevelFatal, 3, err)
        os.Exit(1)
    }
}

//lkj add:
func formatOutput(level int, v ...interface{}) string {
    if runtime.GOOS == "windows" {
        return fmt.Sprintf("%s: %s",
            LevelName(level),
            fmt.Sprint(v...),
        )
    } else {
        return fmt.Sprintf("%s%s: %s%s",
            levelColorPrefix(level),
            LevelName(level),
            fmt.Sprint(v...),
            levelColorSuffix())
    }
}

func (l *Logger) Output(level, calldepth int, v ...interface{}) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if level >= l.level {

        outString := formatOutput(level, v)

        //if runtime.GOOS == "windows"{
        //    h := colorLevelStart_win(level)
        //    defer colorLevelEnd_win(h)
        //    return l.logger.Output(calldepth,
        //       outString,
        //    )
        //} else {
        return l.logger.Output(calldepth,
            outString,
        )
        //}
    }
    return nil
}

// todo: 2020/03/12 error 存入文档
func (l *Logger) Outputf(level, calldepth int, format string, v ...interface{}) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if level >= l.level {
        //lkj modify:
        //		return l.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), fmt.Sprintf(format, v...)))
        //-->
        return l.logger.Output(calldepth,
            formatOutput(level, fmt.Sprintf(format, v...)),
        )
        //]]
    }
    return nil
}

func (l *Logger) Outputln(level, calldepth int, v ...interface{}) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    if level >= l.level {
        s := fmt.Sprintln(v...)
        s = s[:len(s)-1]
        //		return l.logger.Output(calldepth, fmt.Sprintf("%s: %s", LevelName(level), s))
        //-->
        return l.logger.Output(calldepth,
            formatOutput(level, s),
        )
        //]]
    }
    return nil
}

func (l *Logger) Debug(v ...interface{}) {
    l.Output(LevelDebug, 3, v...)
}

func (l *Logger) Info(v ...interface{}) {
    l.Output(LevelInfo, 3, v...)
}

func (l *Logger) Warning(v ...interface{}) {
    l.Output(LevelWarning, 3, v...)
}

func (l *Logger) Error(v ...interface{}) {
    l.Output(LevelError, 3, v...)
}

func (l *Logger) Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    l.Output(LevelPanic, 3, s)
    panic(s)
}

func (l *Logger) Fatal(v ...interface{}) {
    l.Output(LevelFatal, 3, v...)
    os.Exit(1)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
    l.Outputf(LevelDebug, 3, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.Outputf(LevelInfo, 3, format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
    l.Outputf(LevelWarning, 3, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.Outputf(LevelError, 3, format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    l.Outputf(LevelPanic, 3, "%s", s)
    panic(s)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.Outputf(LevelFatal, 3, format, v...)
    os.Exit(1)
}

func (l *Logger) Debugln(v ...interface{}) {
    l.Outputln(LevelDebug, 3, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
    l.Outputln(LevelInfo, 3, v...)
}

func (l *Logger) Warningln(v ...interface{}) {
    l.Outputln(LevelWarning, 3, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
    l.Outputln(LevelError, 3, v...)
}

func (l *Logger) Panicln(v ...interface{}) {
    s := fmt.Sprintln(v...)
    s = s[:len(s)-1]
    l.Outputln(LevelPanic, 3, s)
    panic(s)
}

func (l *Logger) Fatalln(v ...interface{}) {
    l.Outputln(LevelFatal, 3, v...)
    os.Exit(1)
}

//lkj: set level to DEBUG
var std = New(os.Stderr, "", log.LstdFlags|log.Lshortfile, LevelDebug)

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

func DebugDepth(callDepth int, v ...interface{}) {
    std.Output(LevelDebug, callDepth, v...)
}

func Info(v ...interface{}) {
    std.Output(LevelInfo, 3, v...)
}

func InfoDepth(callDepth int, v ...interface{}) {
    std.Output(LevelInfo, callDepth, v...)
}

func Warning(v ...interface{}) {
    std.Output(LevelWarning, 3, v...)
}

func WarningDepth(callDepth int, v ...interface{}) {
    std.Output(LevelWarning, callDepth, v...)
}

func Error(v ...interface{}) {
    std.Output(LevelError, 3, v...)
}

func ErrorDepth(callDepth int, v ...interface{}) {
    std.Output(LevelError, callDepth, v...)
}

//]]

func Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    std.Output(LevelPanic, 3, s)
    panic(s)
}

func PanicDepth(callDepth int, v ...interface{}) {
    s := fmt.Sprint(v...)
    std.Output(LevelPanic, callDepth, s)
    panic(s)
}

////////////////////////////////////////////////////////////////////////////////////////

func Fatal(v ...interface{}) {
    std.Output(LevelFatal, 3, v...)
    os.Exit(1)
}

func FatalDepth(calldepth int, v ...interface{}) {
    std.Output(LevelFatal, calldepth, v...)
    os.Exit(1)
}

func Debugf(format string, v ...interface{}) {
    std.Outputf(LevelDebug, 3, format, v...)
}

func DebugDepthf(calldepth int, format string, v ...interface{}) {
    std.Outputf(LevelDebug, calldepth, format, v...)
}

func Infof(format string, v ...interface{}) {
    std.Outputf(LevelInfo, 3, format, v...)
}

func InfoDepthf(calldepth int, format string, v ...interface{}) {
    std.Outputf(LevelInfo, calldepth, format, v...)
}

func Warningf(format string, v ...interface{}) {
    std.Outputf(LevelWarning, 3, format, v...)
}

func WarningDepthf(calldepth int, format string, v ...interface{}) {
    std.Outputf(LevelWarning, calldepth, format, v...)
}

func Errorf(format string, v ...interface{}) {
    std.Outputf(LevelError, 3, format, v...)
}

func ErrorDepthf(calldepth int, format string, v ...interface{}) {
    std.Outputf(LevelError, calldepth, format, v...)
}

func Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    std.Outputf(LevelPanic, 3, "%s", s)
    panic(s)
}

func PanicDepthf(calldepth int, format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    std.Outputf(LevelPanic, calldepth, "%s", s)
    panic(s)
}

func Fatalf(format string, v ...interface{}) {
    std.Outputf(LevelFatal, 3, format, v...)
    os.Exit(1)
}

func FatalDepthf(calldepth int, format string, v ...interface{}) {
    std.Outputf(LevelFatal, calldepth, format, v...)
    os.Exit(1)
}

////////////////////////////////////////////////////////////////////////////////////////

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
    s = s[:len(s)-1]
    std.Outputln(LevelPanic, 3, s)
    panic(s)
}

func Fatalln(v ...interface{}) {
    std.Outputln(LevelFatal, 3, v...)
    os.Exit(1)
}

////////////////////////////////////////////////////////////////////////////////////////

func Assert(cond bool, a ...interface{}) {
    if !cond {
        PanicDepth(4, a...)
    }
}

func Assertf(cond bool, format string, a ...interface{}) {
    if !cond {
        PanicDepthf(4, format, a...)
    }
}

func AssertDepth(callDepth int, cond bool, a ...interface{}) {
    if !cond {
        PanicDepth(callDepth, a...)
    }
}

func AssertDepthf(callDepth int, cond bool, format string, a ...interface{}) {
    if !cond {
        PanicDepthf(callDepth, format, a...)
    }
}
