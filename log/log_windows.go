package log

import (
    "syscall"
)

var (
    kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
    proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
    CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
)

func init_win() {

}

const (
    colorBlack int = iota
    colorBlue
    colorGreen
    colorCyan
    colorRed
    colorPurple
    colorYellow
    colorLightGray
    colorGray
    colorLightBlue
    colorLightGreen
    colorLightCyan
    colorLightRed
    colorLightPurple
    colorLightYellow
    colorWhite
)

var (
    levels_win = map[int]int{
        LevelDebug:   colorBlue,
        LevelInfo:    colorGreen,
        LevelWarning: colorBlue,
        LevelError:   colorRed,
        LevelPanic:   colorRed,
        LevelFatal:   colorRed,
    }
)

//// 输出有颜色的字体
//func ColorPrint(s string, i int) {
//	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
//	print(s)
//	CloseHandle.Call(handle)
//}

func colorLevelStart_win(level int) uintptr {
    color := levels_win[level]
    handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(color))
    return handle
}

func colorLevelEnd_win(handle uintptr) {
    CloseHandle.Call(handle)
}
