package kitty

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Serialize PNG image from io.Reader into Kitty terminal in-band format.
func kittyCopyPNGInline(out io.Writer, in io.Reader, opts KittyImgOpts) error {

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
