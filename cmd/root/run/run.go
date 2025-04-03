package run

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"os"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/logger"
	"qqlx/cmd"
	"qqlx/pkg/jwt"
)

var Cmd = &cobra.Command{
	Use:   "run",
	Long:  "start the go web framework",
	Short: "start the go web framework",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed(constant.FlagConfigPath) {
			envConfigPath := os.Getenv(constant.ConfigEnv)
			if envConfigPath != "" {
				err := cmd.Flags().Set(constant.FlagConfigPath, envConfigPath)
				if err != nil {
					log.Fatalf("set config file path from env %s faild: %v", envConfigPath, err)
					return
				}
			}
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			cf  string
			err error
		)
		cf, err = cmd.Flags().GetString(constant.FlagConfigPath)
		if err != nil {
			return err
		}
		if cf == "" {
			return errors.New("config file path is empty")
		}
		return runApp(cf)
	},
}

func runApp(configPath string) error {
	err := conf.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("load config file %s faild: %v", configPath, err)
	}
	logger.InitLogger()
	err = jwt.InitConf()
	if err != nil {
		return err
	}
	ctx := context.Background()
	application, cleanup, err := cmd.InitApplication(ctx)
	defer func() {
		_ = zap.S().Sync()
		cleanup()
	}()
	if err != nil {
		return err
	}
	if err = application.Run(ctx); err != nil {
		return err
	}
	return nil
}
