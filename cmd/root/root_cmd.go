package root

import (
	"log"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/cmd/root/init_data"
	"qqlx/cmd/root/run"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     conf.GetProjectName(),
	Long:    `go web framework`,
	Version: constant.ServerVersion,
}

func init() {
	// 添加全局标志
	rootCmd.PersistentFlags().StringP(constant.FlagConfigPath, "C", "./config.yaml", "config file path")
	rootCmd.AddCommand(run.Cmd, init_data.InitCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
