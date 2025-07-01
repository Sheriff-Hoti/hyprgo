package icat

import (
	"github.com/Sheriff-Hoti/hyprgo/consts"
	"github.com/Sheriff-Hoti/hyprgo/pkg"
)

func RenderImages(filenames []string) {
	for idx, filename := range filenames {

		pkg.ICatCmdBuilder(
			pkg.WithScaleUp(),
			pkg.WithStdIn(false),
			pkg.WithImageID(idx+2),
			pkg.WithPlace(
				pkg.Place{
					Width:  consts.ICAT_IMAGE_WIDTH,
					Height: consts.ICAT_IMAGE_HEIGHT,
					Top:    consts.ICAT_IMAGE_TOP_OFFSET + ((idx / consts.CELL_COLS) * 8),
					Left:   consts.ICAT_IMAGE_LEFT_OFFSET + ((idx % consts.CELL_COLS) * (consts.ICAT_IMAGE_WIDTH + 3)),
				},
			),
			pkg.WithWallpaperPath(filename),
		)
	}
}
