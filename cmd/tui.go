package cmd

import (
	"errors"
	"fmt"

	"github.com/Sheriff-Hoti/hyprgo/consts"
	"github.com/Sheriff-Hoti/hyprgo/pkg"
	"github.com/Sheriff-Hoti/hyprgo/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.Flags().StringP("dir", "d", "", "specify the wallpaper directory")
	tuiCmd.Flags().StringP("backend", "b", "", "specify the wallpaper backend")
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "open the tui app",
	Long:  `Well well well`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//prolly here get the config dir or put the default config in
		dir, dir_err := cmd.Flags().GetString("dir")
		backend, backend_err := cmd.Flags().GetString("backend")
		config, config_err := rootCmd.PersistentFlags().GetString("config")

		if flag_err := errors.Join(dir_err, backend_err, config_err); flag_err != nil {
			return flag_err
		}

		if config == "" {
			config = pkg.GetDefaultConfigPath()
		}

		//TODO need to change this so the config dir to be configurable
		kvpairmap, kvpairmap_err := pkg.ReadConfigFile(config)

		if kvpairmap_err != nil {
			return kvpairmap_err
		}

		if backend != "" {
			kvpairmap["backend"] = backend
		}

		if dir != "" {
			kvpairmap["wallpaper_dir"] = dir
		}

		wp_backend := pkg.InitBackend(kvpairmap)

		//TODO make this more dynamic
		filenames, filenames_error := pkg.GetWallpapers(kvpairmap["wallpaper_dir"])

		if filenames_error != nil {
			return filenames_error
		}

		if len(filenames) == 0 {
			//If no png or jpeg or jpg return error
			return errors.New("no wallpapers in this directory")
		}

		RenderImages(filenames)

		fmt.Print("\033[H")

		p := tea.NewProgram(tui.InitialModel(filenames, 0, func(t int) {
			wp_backend.SetImage(filenames[t])
			pkg.DataContent(pkg.DataAction{
				Mode: pkg.Write,
				Data: &pkg.Data{
					Current_wallpaper: filenames[t],
				},
			})

		}))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			return err
		}

		return nil
	},
}

func RenderImages(filenames []string) {
	for idx, filename := range filenames {
		pkg.IcatCmdHalder(pkg.ICatOptions{
			Place: pkg.Place{
				Width:  consts.ICAT_IMAGE_WIDTH,
				Height: consts.ICAT_IMAGE_HEIGHT,
				Top:    consts.ICAT_IMAGE_TOP_OFFSET + ((idx / consts.CELL_COLS) * 8),
				Left:   consts.ICAT_IMAGE_LEFT_OFFSET + ((idx % consts.CELL_COLS) * (consts.ICAT_IMAGE_WIDTH + 3)),
			},
			Extra_args:     []string{"--z-index=--1"},
			Scale_up:       true,
			Wallpaper_path: filename,
		})
	}
}
