package color

import (
	"fmt"
	"runtime"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33;1m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var Black = "\033[30m"
var BrightBlack = "\033[30;1m"
var BrightRed = "\033[31;1m"
var BrightGreen = "\033[32;1m"
var BrightYellow = "\033[33;1m"
var BrightBlue = "\033[34;1m"
var BrightPurple = "\033[35;1m"
var BrightCyan = "\033[36;1m"
var BrightGray = "\033[37;1m"

var BgBlack = "\033[40m"
var BgRed = "\033[41m"
var BgGreen = "\033[42m"
var BgYellow = "\033[43m"
var BgBlue = "\033[44m"
var BgPurple = "\033[45m"
var BgCyan = "\033[46m"
var BgGray = "\033[47m"

var BgBrightBlack = "\033[100m"
var BgBrightRed = "\033[101m"
var BgBrightGreen = "\033[102m"
var BgBrightYellow = "\033[103m"
var BgBrightBlue = "\033[104m"
var BgBrightPurple = "\033[105m"
var BgBrightCyan = "\033[106m"
var BgBrightGray = "\033[107m"

var Bold = "\033[1m"
var Dim = "\033[2m"
var Italic = "\033[3m"
var Underline = "\033[4m"
var Blink = "\033[5m"
var Inverse = "\033[7m"
var Hidden = "\033[8m"
var Strikethrough = "\033[9m"

var E_ORANGE = RGBColor(228, 129, 11)
var E_YELLOW = RGBColor(247, 247, 17)
var E_GREEN = RGBColor(167, 237, 135)
var E_BLUE = Cyan
var E_MUTE = RGBColor(77, 77, 55)

// color needs to be between 0-255
var CustomColor = func(color int) string {
	if color < 0 || color > 256 {
		color = 99
	}
	return fmt.Sprintf("\033[38;5;%dm", color)
}

var CustomBgColor = func(color int) string {
	if color < 0 || color > 256 {
		color = 99
	}
	return fmt.Sprintf("\033[48;5;%dm", color)
}

var RGBColor = func(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}
var RGBBgColor = func(r, g, b int) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}
