package tui

import (
	"fmt"
	"io"
	"math/rand/v2"

	"github.com/Sheriff-Hoti/hyprgo/kitty"
)

type Page struct {
	cells         []kitty.Cell
	grid_opts     kitty.KittyGridOpts
	order         int
	selected_cell int
	initialized   bool
}

func NewPage(order int, abs_filenames []string, grid_opts kitty.KittyGridOpts) *Page {
	cells := make([]kitty.Cell, 0, len(abs_filenames))
	fnames_idx := 0
outer_loop:
	for row_idx := range grid_opts.Rows {
		for col_idx := range grid_opts.Cols {

			if fnames_idx >= len(abs_filenames) {
				break outer_loop
			}
			row_cell := (row_idx * grid_opts.ImgHeight) + grid_opts.TopSpacing + (grid_opts.RowsSpacing * row_idx)
			col_cell := (col_idx * grid_opts.ImgWidth) + grid_opts.LeftSpacing + (grid_opts.ColsSpacing * col_idx)

			cells = append(cells, kitty.Cell{
				Filename: abs_filenames[fnames_idx],
				Width:    uint32(grid_opts.ImgWidth),
				Height:   uint32(grid_opts.ImgHeight),
				RowCell:  uint32(row_cell),
				ColCell:  uint32(col_cell),
				Id:       uint32(rand.IntN(100)),
			})

			fnames_idx++

		}

	}

	return &Page{
		cells:     cells,
		grid_opts: grid_opts,
		order:     order,
	}
}

func (p *Page) RenderImages(out io.Writer) error {
	for _, cell := range p.cells {
		_, erro := fmt.Fprintf(out, "\x1b[%d;%dH", cell.RowCell, cell.ColCell)
		if erro != nil {
			return erro
		}

		err := cell.RenderImage(out, kitty.KittyImgOpts{
			DstCols:     uint32(p.grid_opts.ImgWidth),
			DstRows:     uint32(p.grid_opts.ImgHeight),
			CellOffsetX: 0,
			CellOffsetY: 0,
			//TODO: fix it
			ImageId:     cell.Id,
			PlacementId: cell.Id,
		})

		if err != nil {
			return err
		}
	}
	p.initialized = true
	return nil
}

func (p *Page) Hide(out io.Writer) error {
	for _, cell := range p.cells {
		cell.Hide(out, kitty.KittyImgOpts{
			ImageId:     cell.Id,
			PlacementId: cell.Id,
		})
	}
	return nil
}

func (p *Page) Show(out io.Writer) error {

	for _, cell := range p.cells {
		fmt.Fprintf(out, "\x1b[%d;%dH", cell.RowCell, cell.ColCell)

		cell.Show(out, kitty.KittyImgOpts{
			DstCols:     uint32(p.grid_opts.ImgWidth),
			DstRows:     uint32(p.grid_opts.ImgHeight),
			CellOffsetX: 0,
			CellOffsetY: 0,
			//TODO: fix it
			ImageId:     cell.Id,
			PlacementId: cell.Id,
		})
	}

	return nil
}
