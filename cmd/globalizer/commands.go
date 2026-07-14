package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("globalizer v%s\n", version)
		fmt.Println("OpenSource Globalizer AI — 面向开源项目的 AI 国际化助手")
	},
}

var languagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "列出支持的语言",
	Run: func(cmd *cobra.Command, args []string) {
		langs := []struct {
			Code, Name, Native string
		}{
			{"en", "English", "English"},
			{"zh-CN", "Chinese (Simplified)", "简体中文"},
			{"ja", "Japanese", "日本語"},
			{"ko", "Korean", "한국어"},
			{"es", "Spanish", "Español"},
			{"fr", "French", "Français"},
			{"de", "German", "Deutsch"},
			{"pt-BR", "Portuguese (Brazil)", "Português (Brasil)"},
			{"ru", "Russian", "Русский"},
			{"ar", "Arabic", "العربية"},
		}
		fmt.Println("支持的语言:")
		for _, l := range langs {
			fmt.Printf("  %-6s  %-25s %s\n", l.Code, l.Name, l.Native)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(languagesCmd)
}
