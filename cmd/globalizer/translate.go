package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ytc301/opensource-globalizer/internal/ai"
	"github.com/ytc301/opensource-globalizer/internal/store"
	"github.com/ytc301/opensource-globalizer/internal/translator"
	"github.com/ytc301/opensource-globalizer/pkg/config"
	"go.uber.org/zap"
)

var (
	targetLangs string
	outputDir   string
	model       string
	dryRun      bool
	sourceLang  string
	useMock     bool

	translateCmd = &cobra.Command{
		Use:   "translate <file>",
		Short: "翻译 Markdown 文件为多语言版本",
		Long: `将 README.md 等 Markdown 文件翻译为指定的多语言版本。

示例:
  globalizer translate README.md --lang zh-CN,ja,ko
  globalizer translate README.md --lang zh-CN --output i18n/
  globalizer translate README.md --lang zh-CN --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: runTranslate,
	}
)

func init() {
	translateCmd.Flags().StringVarP(&targetLangs, "lang", "l", "zh-CN", "目标语言，逗号分隔")
	translateCmd.Flags().StringVarP(&outputDir, "output", "o", "docs", "输出目录")
	translateCmd.Flags().StringVarP(&model, "model", "m", "gpt-4o", "OpenAI 模型名称")
	translateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览模式，不写入文件")
	translateCmd.Flags().StringVarP(&sourceLang, "source", "s", "", "源语言 (留空自动检测)")
	translateCmd.Flags().BoolVar(&useMock, "mock", false, "使用 Mock Provider (无需 API Key，仅供测试)")

	rootCmd.AddCommand(translateCmd)
}

func runTranslate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// 加载配置
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 解析目标语言
	langList := strings.Split(targetLangs, ",")
	for i := range langList {
		langList[i] = strings.TrimSpace(langList[i])
	}

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// 创建 AI Provider
	var provider ai.Provider
	if useMock {
		provider = ai.NewMockProvider()
		fmt.Println("🧪 使用 Mock Provider (测试模式)")
	} else {
		apiKey := cfg.EffectiveAPIKey()
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			return fmt.Errorf("请设置 OPENAI_API_KEY 环境变量或在配置文件中指定")
		}
		provider = ai.NewOpenAI(ai.OpenAIConfig{
			APIKey:  apiKey,
			BaseURL: cfg.OpenAI.BaseURL,
		})
	}

	// 创建 Store (可选，缓存翻译)
	var st *store.Store
	if cfg.DBPath != "" {
		st, err = store.New(cfg.DBPath)
		if err != nil {
			logger.Warn("无法创建存储，跳过翻译缓存", zap.Error(err))
			st = nil
		}
	}
	if st != nil {
		defer st.Close()
	}

	// 创建翻译服务
	svc := translator.NewService(provider, st, logger)

	fmt.Printf("📖 源文件: %s\n", filePath)
	fmt.Printf("🌍 目标语言: %s\n", strings.Join(langList, ", "))
	fmt.Printf("🤖 模型: %s\n", model)
	if dryRun {
		fmt.Println("🔍 预览模式 (不写入文件)")
	}

	// 执行翻译
	results, err := svc.TranslateFile(cmd.Context(), string(content), langList, model)
	if err != nil {
		return fmt.Errorf("翻译失败: %w", err)
	}

	// 输出结果
	if dryRun {
		for lang, translated := range results {
			limit := len(translated)
			if limit > 500 {
				limit = 500
			}
			fmt.Printf("\n--- %s ---\n%s\n", lang, translated[:limit])
		}
		return nil
	}

	// 写入文件
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for lang, translated := range results {
		outPath := fmt.Sprintf("%s/README.%s.md", outputDir, lang)
		if err := os.WriteFile(outPath, []byte(translated), 0644); err != nil {
			return fmt.Errorf("写入 %s 失败: %w", outPath, err)
		}
		fmt.Printf("  ✅ %s\n", outPath)
	}

	fmt.Printf("\n✨ 翻译完成！共 %d 个语言版本\n", len(results))
	return nil
}
