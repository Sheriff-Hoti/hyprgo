package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

// handle here kitten icat callings

var (
	base_cmd = "kitten"
	icat_cmd = "icat"
)

//kitten icat --stdin=no --align=left --place=20x20@1x3 ./img/night-sky.jpg
//&& kitten icat --stdin=no --align=left --place=20x20@1x10 ./img/test.jpg

type ICatOptions struct {
	stdin    bool
	scale_up bool
	place    struct {
		width  int
		height int
		left   int
		top    int
	}
	extra_args     []string
	wallpaper_path string
}

// https://github.com/5hubham5ingh/WallRizz/blob/main/src/userInterface.js
func IcatCmdHalder(options ICatOptions) {
	stdin := "--stdin=no"
	scale_up := ""
	place := fmt.Sprintf("--place=%vx%v@%vx%v", options.place.width, options.place.height, options.place.left, options.place.top)
	if options.stdin {
		stdin = "--stdin=yes"
	}
	if options.scale_up {
		scale_up = "--scale-up"
	}

	fullArgs := append([]string{icat_cmd, stdin, scale_up, place}, options.extra_args...)

	cmd := exec.Command(base_cmd, append(fullArgs, options.wallpaper_path)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
}

// drawUI() {
//     if (!this.wallpapers) return;
//     print(clearTerminal);
//     // Draw wallpapers
//     this.wallpapers.forEach((wallpaper, i) => {
//       const wallpaperDir = `${this.wallpapersDir}/${wallpaper.uniqueId}`;
//       const [x, y] = i < this.xy.length
//         ? this.xy[i]
//         : this.xy[i % this.xy.length];
//       const coordinates = `${this.imageWidth}x${this.imageHeight}@${x}x${y}`;
//       // print(cursorMove(x, y));
//       // OS.exec([
//       //   "timg",
//       //   "-U",
//       //   "-W",
//       //   "--clear",
//       //   "-pk",
//       //   `-g${this.imageWidth}x${this.imageHeight}`,
//       //   wallpaperDir,
//       // ]);
//       OS.exec([
//         "kitten",
//         "icat",
//         "--stdin=no",
//         "--scale-up",
//         "--place",
//         coordinates,
//         wallpaperDir,
//       ]);
//     });

//     this.drawContainerBorder(this.xy[this.selection]);
//   }
