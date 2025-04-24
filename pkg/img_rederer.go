package pkg

import (
	"fmt"
	"os"

	"github.com/BourgeoisBear/rasterm"
)

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

func Store_temp() {
	fIn, _, err := getFile("./night-sky.jpg")
	if err != nil {
		println(err)
	}

	fmt.Println("\nKitty PNG Inline")
	//this totally works
	//TODO need to provide the geometry the starting points and the width and height of the images
	eI := rasterm.KittyCopyPNGInline(os.Stdout, fIn, rasterm.KittyImgOpts{})
	rasterm.KittyCopyPNGInline(os.Stdout, fIn, rasterm.KittyImgOpts{})
	rasterm.KittyCopyPNGInline(os.Stdout, fIn, rasterm.KittyImgOpts{})
	rasterm.KittyCopyPNGInline(os.Stdout, fIn, rasterm.KittyImgOpts{})
	rasterm.KittyCopyPNGInline(os.Stdout, fIn, rasterm.KittyImgOpts{})

	print(eI)
}
