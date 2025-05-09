package cmd

import (
	"github.com/Sheriff-Hoti/hyprgo/pkg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "open the tui app",
	Long:  `Well well well`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data_content, data_err := pkg.DataContent(pkg.DataAction{
			Mode: pkg.Read,
			Data: nil,
		})

		if data_err != nil {
			return data_err
		}

		kvpairmap, config_err := pkg.ReadConfigFile(nil)

		if config_err != nil {
			return config_err
		}

		wp_backend := pkg.InitBackend(kvpairmap)

		err := wp_backend.SetImage(data_content.Current_wallpaper)

		if err != nil {
			return err
		}

		return nil
	},
}
