package cmd

import (
	"errors"

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

		config, config_err := rootCmd.PersistentFlags().GetString("config")

		data_content, data_err := pkg.DataContent(pkg.DataAction{
			Mode: pkg.Read,
			Data: nil,
		})

		if flag_err := errors.Join(data_err, config_err); flag_err != nil {
			return flag_err
		}

		if !cmd.Flags().Changed("config") {
			config = pkg.GetDefaultConfigPath()
		}

		kvpairmap, config_err := pkg.ReadConfigFile(config)

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
