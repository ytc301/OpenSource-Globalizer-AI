package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cfgFile string
	logger  *zap.Logger
	rootCmd = &cobra.Command{
		Use:   "globalizer",
		Short: "🌍 OpenSource Globalizer AI — 面向开源项目的 AI 国际化助手",
		Long: `OpenSource Globalizer AI helps open-source maintainers build truly global
communities by automating documentation localization and community support.

Commands:
  translate   Translate Markdown files to multiple languages
  languages   List supported languages
  serve       Start HTTP API server
  version     Show version information`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: .globalizer.yaml)")
}

func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	initLogger()
	defer logger.Sync()

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("命令执行失败", zap.Error(err))
	}
}
