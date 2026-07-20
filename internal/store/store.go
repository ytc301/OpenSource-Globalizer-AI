package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Store 基于 GORM + SQLite 的数据访问层。
type Store struct {
	db *gorm.DB
}

// Translation 翻译缓存记录。
type Translation struct {
	ID         uint   `gorm:"primaryKey"`
	SourceHash string `gorm:"uniqueIndex:idx_source_target;not null"`
	TargetLang string `gorm:"uniqueIndex:idx_source_target;not null"`
	SourceText string `gorm:"not null"`
	Translated string `gorm:"not null"`
	Model      string `gorm:"not null"`
	TokensUsed int
}

// New 创建 Store 并自动迁移。
func New(dbPath string) (*Store, error) {
	if strings.HasPrefix(dbPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录: %w", err)
		}
		dbPath = filepath.Join(home, dbPath[2:])
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库: %w", err)
	}

	if err := db.AutoMigrate(&Translation{}); err != nil {
		return nil, fmt.Errorf("数据库迁移: %w", err)
	}

	return &Store{db: db}, nil
}

// Close 关闭数据库连接。
func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetCached 查找翻译缓存。
func (s *Store) GetCached(sourceHash, targetLang string) (*Translation, error) {
	var t Translation
	result := s.db.Where("source_hash = ? AND target_lang = ?", sourceHash, targetLang).First(&t)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询缓存: %w", result.Error)
	}
	return &t, nil
}

// PutCache 写入或更新翻译缓存。
func (s *Store) PutCache(t *Translation) error {
	return s.db.Where("source_hash = ? AND target_lang = ?", t.SourceHash, t.TargetLang).
		Assign(map[string]interface{}{
			"source_text": t.SourceText,
			"translated":  t.Translated,
			"model":       t.Model,
			"tokens_used": t.TokensUsed,
		}).
		FirstOrCreate(t).Error
}

// DB 返回底层 GORM 实例。
func (s *Store) DB() *gorm.DB {
	return s.db
}
