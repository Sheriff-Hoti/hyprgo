package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/consts"
	"github.com/Sheriff-Hoti/hyprgo/pkg"
	"github.com/Sheriff-Hoti/hyprgo/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyprgo",
	Short: "tui wallpaper picker",
	Long: `tui wallpaper picker:

well well well, the design is inspired by wallrizz".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, dir_err := cmd.Flags().GetString("dir")
		backend, backend_err := cmd.Flags().GetString("backend")
		config, config_err := cmd.PersistentFlags().GetString("config")

		if flag_err := errors.Join(dir_err, backend_err, config_err); flag_err != nil {
			return flag_err
		}

		if config == "" {
			config = pkg.GetDefaultConfigPath()
		}

		//TODO need to change this so the config dir to be configurable
		cc, cc_err := pkg.ReadConfigFile(config)

		if cc_err != nil {
			return cc_err
		}

		if cmd.Flags().Changed("backend") {
			cc.Backend = backend
		}

		if cmd.Flags().Changed("dir") {
			cc.Wallpaper_dir = dir
		}

		wp_backend := pkg.InitBackend(cc)

		//TODO make this more dynamic
		filenames, filenames_error := pkg.GetWallpapers(cc.Wallpaper_dir)

		if filenames_error != nil {
			return filenames_error
		}

		// if len(filenames) == 0 {
		// 	//If no png or jpeg or jpg return error
		// 	return errors.New("no wallpapers in this directory")
		// }

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "specify the config file")
	rootCmd.Flags().StringP("dir", "d", "", "specify the wallpaper directory")
	rootCmd.Flags().StringP("backend", "b", "", "specify the wallpaper backend")
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
