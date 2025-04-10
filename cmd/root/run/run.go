package run

import (
	"context"
	"errors"
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cf  string
			err error
		)
		cf, err = cmd.Flags().GetString(constant.FlagConfigPath)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if cf == "" {
			log.Fatal("config file path is empty")
		}
		runApp(cf)
	},
}

func runApp(configPath string) {
	err := conf.LoadConfig(configPath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("load config file %s faild: %v", configPath, err))
	}
	logger.InitLogger()
	err = jwt.InitConf()
	err = errors.New("testTEst")
	if err != nil {
		zap.S().Fatal(err)
	}
	ctx := context.Background()
	application, cleanup, err := cmd.InitApplication(ctx)
	defer func() {
		_ = zap.S().Sync()
		cleanup()
	}()
	if err != nil {
		zap.S().Fatal(err)
	}
	if err = application.Run(ctx); err != nil {
		zap.S().Fatal(err)
	}
}
