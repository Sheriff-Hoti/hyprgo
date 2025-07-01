package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/cmd"
	kitty "github.com/Sheriff-Hoti/hyprgo/kitty"
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

	// cmd.Execute()
	if fmtName == "jpg" {

		fmt.Println("Kitty PNG Local File")
		kitty.KittyWritePNGLocal(os.Stdout, fpath, kitty.KittyImgOpts{})
		// fmt.Println("\nKitty PNG Inline")
		// eI := kitty.KittyCopyPNGInline(os.Stdout, fIn, kitty.KittyImgOpts{})
		// err = errors.Join(eI, eF)

	} else {

		err = kitty.KittyWriteImage(os.Stdout, iImg, kitty.KittyImgOpts{})
	}
}

func main() {
	cmd.Execute()

}
