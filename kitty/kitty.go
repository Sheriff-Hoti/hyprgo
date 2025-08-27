package kitty

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html for more details.

type KittyGridOpts struct {
	Cols        int // display width in terminal columns
	Rows        int // display height in terminal rows
	RowsSpacing int // spacing between rows in terminal rows
	ColsSpacing int // spacing between columns in terminal columns
	ImgWidth    int // image width in pixels
	ImgHeight   int // image height in pixels
}

const (
	KITTY_IMG_HDR = "\x1b_G"
	KITTY_IMG_FTR = "\x1b\\"
)

type KittyImgOpts struct {
	SrcX        uint32 // x=
	SrcY        uint32 // y=
	SrcWidth    uint32 // w=
	SrcHeight   uint32 // h=
	CellOffsetX uint32 // X= (pixel x-offset inside terminal cell)
	CellOffsetY uint32 // Y= (pixel y-offset inside terminal cell)
	DstCols     uint32 // c= (display width in terminal columns)
	DstRows     uint32 // r= (display height in terminal rows)
	ZIndex      int32  // z=
	ImageId     uint32 // i=
	ImageNo     uint32 // I=
	PlacementId uint32 // p=
}

func (o KittyImgOpts) ToHeader(opts ...string) string {

	type fldmap struct {
		pv   *uint32
		code rune
	}
	sFld := []fldmap{
		{&o.SrcX, 'x'},
		{&o.SrcY, 'y'},
		{&o.SrcWidth, 'w'},
		{&o.SrcHeight, 'h'},
		{&o.CellOffsetX, 'X'},
		{&o.CellOffsetY, 'Y'},
		{&o.DstCols, 'c'},
		{&o.DstRows, 'r'},
		{&o.ImageId, 'i'},
		{&o.ImageNo, 'I'},
		{&o.PlacementId, 'p'},
	}

	for _, f := range sFld {
		if *f.pv != 0 {
			opts = append(opts, fmt.Sprintf("%c=%d", f.code, *f.pv))
		}
	}

	if o.ZIndex != 0 {
		opts = append(opts, fmt.Sprintf("z=%d", o.ZIndex))
	}

	return KITTY_IMG_HDR + strings.Join(opts, ",") + ";"
}

// checks if terminal supports kitty image protocols
func IsKittyCapable() bool {

	// TODO: more rigorous check
	V := GetEnvIdentifiers()
	return (len(V["KITTY_WINDOW_ID"]) > 0) || (V["TERM_PROGRAM"] == "wezterm") || (V["TERM_PROGRAM"] == "ghostty")
}

// Display local PNG file
// - pngFileName must be directly accesssible from Kitty instance
// - pngFileName must be an absolute path
func KittyWritePNGLocal(out io.Writer, pngFileName string, opts KittyImgOpts) error {

	_, e := fmt.Fprint(out, opts.ToHeader("a=T", "f=100", "t=f"))
	if e != nil {
		return e
	}

	enc64 := base64.NewEncoder(base64.StdEncoding, out)

	_, e = fmt.Fprint(enc64, pngFileName)
	if e != nil {
		return e
	}

	e = enc64.Close()
	if e != nil {
		return e
	}

	_, e = fmt.Fprint(out, KITTY_IMG_FTR)
	return e
}

// Serialize image.Image into Kitty terminal in-band format.
func KittyWriteImage(out io.Writer, iImg image.Image, opts KittyImgOpts) error {

	pBuf := new(bytes.Buffer)
	if E := png.Encode(pBuf, iImg); E != nil {
		return E
	}

	return KittyCopyPNGInline(out, pBuf, opts)
}

// Serialize PNG image from io.Reader into Kitty terminal in-band format.
func KittyCopyPNGInline(out io.Writer, in io.Reader, opts KittyImgOpts) error {

	_, err := fmt.Fprint(out, opts.ToHeader("a=T", "f=100", "t=d", "m=1"), KITTY_IMG_FTR)
	if err != nil {
		return err
	}

	// PIPELINE: PNG (io.Reader) -> B64 -> CHUNKER -> (io.Writer)
	// SEND IN 4K CHUNKS
	cw := KittyChunkWri{
		nChunkSize: 4096,
		iWri:       out,
	}

	enc64 := base64.NewEncoder(base64.StdEncoding, &cw)
	_, err = io.Copy(enc64, in)
	return errors.Join(
		err,
		enc64.Close(),
		cw.Close(),
	)
}

func GetEnvIdentifiers() map[string]string {

	KEYS := []string{"TERM", "TERM_PROGRAM", "LC_TERMINAL", "VIM_TERMINAL", "KITTY_WINDOW_ID"}
	V := make(map[string]string)
	for _, K := range KEYS {
		V[K] = lcaseEnv(K)
	}

	return V
}

func lcaseEnv(k string) string {
	return strings.ToLower(strings.TrimSpace(os.Getenv(k)))
}

func KittyWriteFileGPTGENERATED(out io.Writer, fileName string, opts KittyImgOpts) error {
	// Check absolute path
	if !strings.HasPrefix(fileName, "/") {
		return fmt.Errorf("file path must be absolute: %s", fileName)
	}

	// Check if file exists
	if _, err := os.Stat(fileName); err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	// Build the Kitty header
	header := opts.ToHeader("a=T", "f=100", "t=f")

	// Write header
	if _, err := fmt.Fprint(out, header); err != nil {
		return err
	}

	// Encode the absolute path in base64 (required by Kitty)
	enc64 := base64.NewEncoder(base64.StdEncoding, out)
	if _, err := fmt.Fprint(enc64, fileName); err != nil {
		return err
	}
	if err := enc64.Close(); err != nil {
		return err
	}

	// Write the terminal escape sequence to finish
	if _, err := fmt.Fprint(out, KITTY_IMG_FTR); err != nil {
		return err
	}

	return nil
}

func KittyWriteFiles(out io.Writer, fileNames []string, grid KittyGridOpts) error {
	if grid.Cols <= 0 {
		return fmt.Errorf("grid.Cols must be > 0")
	}

	for idx, fileName := range fileNames {

		col := idx % grid.Cols
		row := idx / grid.Cols

		if grid.Rows > 0 && row >= grid.Rows {
			return nil
		}

		log.Println(col, row, idx)

		fmt.Fprintf(out, "\x1b[%d;%dH", (col*grid.ImgHeight)+(grid.ColsSpacing*(col+1)), (row*grid.ImgWidth)+(grid.RowsSpacing*(row+1)))

		absfile, err := filepath.Abs(fileName)
		if err != nil {
			return err
		}

		if err := KittyWriteFileGPTGENERATED(out, absfile, KittyImgOpts{
			DstCols:     uint32(grid.ImgWidth),  // display width in terminal columns
			DstRows:     uint32(grid.ImgHeight), // display height in terminal rows
			CellOffsetX: 0,                      // anchor at the cell, don't use pixel offsets here
			CellOffsetY: 0,
		}); err != nil {
			return err
		}

	}
	return nil
}

// check this:https://chatgpt.com/share/68ae4268-106c-8007-bcb5-16476f778c24

// start := time.Now()
// 	err := kitty.KittyWriteFiles(os.Stdout, []string{
// 		"./test_assets/img/test0.jpg",
// 		"./test_assets/img/test1.jpg",
// 		"./test_assets/img/test2.jpg",
// 		"./test_assets/img/test3.jpg",
// 		"./test_assets/img/test4.jpg",
// 		"./test_assets/img/test5.png",
// 	}, kitty.KittyGridOpts{
// 		Cols:        3,
// 		Rows:        3,
// 		RowsSpacing: 5,
// 		ColsSpacing: 5,
// 		ImgWidth:    consts.ICAT_IMAGE_WIDTH,
// 		ImgHeight:   consts.ICAT_IMAGE_HEIGHT,
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Fprintf(os.Stderr, "\nRendered images in %v\n", time.Since(start).Seconds())
