package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Sheriff-Hoti/hyprgo/consts"
	"github.com/Sheriff-Hoti/hyprgo/kitty"
)

// import (
//
//	"github.com/Sheriff-Hoti/hyprgo/cmd"
//
// )
func getFile(fpath string) (*os.File, int64, error) {

	pF, E := os.Open(fpath)
	if E != nil {
		return nil, 0, E
	}

	fInf, E := pF.Stat()
	if E != nil {
		pF.Close()
		return nil, 0, E
	}

	return pF, fInf.Size(), nil
}

func test() {
	fpath := "./test_assets/img/test0.jpg"
	fIn, nImgLen, err := getFile(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer fIn.Close()

	imgCfg, fmtName, err := image.DecodeConfig(fIn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fIn.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	iImg, _, err := image.Decode(fIn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fIn.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[FMT: %s, W: %d, H: %d, LEN: %d, IMG: %T]\n", fmtName, imgCfg.Width, imgCfg.Height, nImgLen, iImg)

	//this is the answer
	kitty.KittyWritePNGLocal(os.Stdout, fpath, kitty.KittyImgOpts{})

	// cmd.Execute()
	// if fmtName == "jpg" {
	// 	fmt.Println("Kitty PNG Local File")
	// 	kitty.KittyWritePNGLocal(os.Stdout, fpath, kitty.KittyImgOpts{})
	// 	// fmt.Println("\nKitty PNG Inline")
	// 	eI := kitty.KittyCopyPNGInline(os.Stdout, fIn, kitty.KittyImgOpts{})
	// 	err = errors.Join(eI)

	// } else {

	// 	err = kitty.KittyWriteImage(os.Stdout, iImg, kitty.KittyImgOpts{})
	// }
}

func maintest() {
	start := time.Now()

	out := os.Stdout

	// use values from consts for consistent sizing
	squareCols := consts.ICAT_IMAGE_WIDTH
	squareRows := consts.ICAT_IMAGE_HEIGHT
	spacing := 3

	// List of images
	images := []string{
		"./test_assets/img/test0.jpg",
		"./test_assets/img/test1.jpg",
		"./test_assets/img/test2.jpg",
	}

	// Track horizontal position (in columns)
	// ...existing code...
	for colIndex, imgPath := range images {
		absPath, _ := filepath.Abs(imgPath)

		// compute placement using same logic as icat.RenderImages
		// Left  = ICAT_IMAGE_LEFT_OFFSET + ((idx % CELL_COLS) * (ICAT_IMAGE_WIDTH + spacing))
		// Top   = ICAT_IMAGE_TOP_OFFSET  + ((idx / CELL_COLS) * 8)
		left := consts.ICAT_IMAGE_LEFT_OFFSET + ((colIndex % consts.CELL_COLS) * (squareCols + spacing))
		top := consts.ICAT_IMAGE_TOP_OFFSET + ((colIndex / consts.CELL_COLS) * 8)

		// move cursor to target cell (1-based row/col) so kitty anchors image there
		// CSI {row};{col}H
		fmt.Fprintf(out, "\x1b[%d;%dH", top+1, left+1)

		opts := kitty.KittyImgOpts{
			DstCols:     uint32(squareCols), // display width in terminal columns
			DstRows:     uint32(squareRows), // display height in terminal rows
			CellOffsetX: 0,                  // anchor at the cell, don't use pixel offsets here
			CellOffsetY: 0,
		}

		if err := kitty.KittyWriteFileGPTGENERATED(out, absPath, opts); err != nil {
			fmt.Fprintln(os.Stderr, "kitty write error:", err)
		}
		// time.Sleep(time.Second)
	}

	secondimages := []string{
		"./test_assets/img/test3.jpg",
		"./test_assets/img/test4.jpg",
		"./test_assets/img/test5.png",
	}

	// number of rows used by first batch (in terminal rows)
	numRowsFirst := (len(images) + consts.CELL_COLS - 1) / consts.CELL_COLS

	// vertical gap between batches in terminal rows (adjust as needed)
	vGap := 1

	// base top row for second batch
	baseTop := consts.ICAT_IMAGE_TOP_OFFSET + numRowsFirst*(squareRows+vGap)

	for colIndex, imgPath := range secondimages {
		absPath, _ := filepath.Abs(imgPath)

		left := consts.ICAT_IMAGE_LEFT_OFFSET + ((colIndex % consts.CELL_COLS) * (squareCols + spacing))
		top := baseTop + ((colIndex / consts.CELL_COLS) * (squareRows + vGap))

		// move cursor to target cell (1-based row/col) so kitty anchors image there
		fmt.Fprintf(out, "\x1b[%d;%dH", top+1, left+1)

		opts := kitty.KittyImgOpts{
			DstCols:     uint32(squareCols),
			DstRows:     uint32(squareRows),
			CellOffsetX: 0,
			CellOffsetY: 0,
		}

		if err := kitty.KittyWriteFileGPTGENERATED(out, absPath, opts); err != nil {
			fmt.Fprintln(os.Stderr, "kitty write error:", err)
		}
	}

	// thirdimages := []string{
	// 	"./test_assets/img/test3.jpg",
	// 	"./test_assets/img/test4.jpg",
	// 	"./test_assets/img/test4.jpg",
	// }

	// // number of rows used by first batch (in terminal rows)
	// numRowsSecond := (len(secondimages) + consts.CELL_COLS - 1) / consts.CELL_COLS

	// // vertical gap between batches in terminal rows (adjust as needed)

	// // base top row for second batch
	// baseTop2 := consts.ICAT_IMAGE_TOP_OFFSET + numRowsSecond*(squareRows+vGap)

	// for colIndex, imgPath := range thirdimages {
	// 	absPath, _ := filepath.Abs(imgPath)

	// 	left := consts.ICAT_IMAGE_LEFT_OFFSET + ((colIndex % consts.CELL_COLS) * (squareCols + spacing))
	// 	top := baseTop2 + ((colIndex / consts.CELL_COLS) * (squareRows + vGap))

	// 	// move cursor to target cell (1-based row/col) so kitty anchors image there
	// 	fmt.Fprintf(out, "\x1b[%d;%dH", top+1, left+1)

	// 	opts := kitty.KittyImgOpts{
	// 		DstCols:     uint32(squareCols),
	// 		DstRows:     uint32(squareRows),
	// 		CellOffsetX: 0,
	// 		CellOffsetY: 0,
	// 	}

	// 	if err := kitty.KittyWriteFileGPTGENERATED(out, absPath, opts); err != nil {
	// 		fmt.Fprintln(os.Stderr, "kitty write error:", err)
	// 	}
	// }

	duration := time.Since(start) // measure elapsed time

	fmt.Fprintf(os.Stderr, "\nRendered %d images in %v\n", len(images)+len(secondimages), duration.Seconds())
}

// cmd.Execute()

func main() {
	start := time.Now()
	err := kitty.KittyWriteFiles(os.Stdout, []string{
		"./test_assets/img/test0.jpg",
		"./test_assets/img/test1.jpg",
		"./test_assets/img/test2.jpg",
		"./test_assets/img/test3.jpg",
		"./test_assets/img/test4.jpg",
		"./test_assets/img/test5.png",
	}, kitty.KittyGridOpts{
		Cols:        3,
		Rows:        3,
		RowsSpacing: 5,
		ColsSpacing: 5,
		ImgWidth:    consts.ICAT_IMAGE_WIDTH,
		ImgHeight:   consts.ICAT_IMAGE_HEIGHT,
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "\nRendered images in %v\n", time.Since(start).Seconds())

}
