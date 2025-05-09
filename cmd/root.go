package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyprgo",
	Short: "tui wallpaper picker",
	Long: `tui wallpaper picker:

well well well, the design is inspired by wallrizz".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, config_err := cmd.Flags().GetString("config")

		if config_err != nil {
			return config_err
		}
		fmt.Println(config)
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
}
