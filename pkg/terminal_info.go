package pkg

// import (
// 	"os"
// 	"strings"

// 	"golang.org/x/sys/unix"
// )

// type TerminalInfo struct {
// 	rows    uint16
// 	columns uint16
// 	width   uint16
// 	height  uint16
// }

// func GetTerminalSize() (term_info *TerminalInfo, err error) {
// 	f, f_err := os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666)

// 	if f_err != nil {
// 		return nil, f_err
// 	}

// 	var sz *unix.Winsize
// 	sz, unix_err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)

// 	if unix_err != nil {
// 		return nil, unix_err
// 	}

// 	return &TerminalInfo{
// 		rows:    sz.Row,
// 		columns: sz.Col,
// 		width:   sz.Xpixel,
// 		height:  sz.Ypixel,
// 	}, nil

// }

// func lcaseEnv(k string) string {
// 	//borrowed from https://github.com/BourgeoisBear/rasterm/blob/main/term_misc.go#L183
// 	return strings.ToLower(strings.TrimSpace(os.Getenv(k)))
// }

// func GetEnvIdentifiers() map[string]string {
// 	//borrowed from https://github.com/BourgeoisBear/rasterm/blob/main/term_misc.go#L183

// 	KEYS := []string{"TERM", "TERM_PROGRAM", "LC_TERMINAL", "VIM_TERMINAL", "KITTY_WINDOW_ID"}
// 	V := make(map[string]string)
// 	for _, K := range KEYS {
// 		V[K] = lcaseEnv(K)
// 	}

// 	return V
// }

// // IS kitty Terminal Graphics Protocol supported
// func IsKittyTGPSupported() bool {

// 	//borrowed from https://github.com/BourgeoisBear/rasterm/blob/main/kitty.go#L70

// 	// TODO: more rigorous check
// 	V := GetEnvIdentifiers()
// 	return (len(V["KITTY_WINDOW_ID"]) > 0) || (V["TERM_PROGRAM"] == "wezterm") || (V["TERM_PROGRAM"] == "ghostty")
// }

// //TODO search this https://sw.kovidgoyal.net/kitty/graphics-protocol/#relative-placements
