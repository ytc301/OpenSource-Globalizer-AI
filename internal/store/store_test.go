package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	st, err := New(dbPath)
	if err != nil {
		t.Fatalf("创建测试 Store 失败: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

func TestNewStore_AutoMigrate(t *testing.T) {
	st := newTestStore(t)
	if st.db == nil {
		t.Fatal("数据库连接为空")
	}
	// 验证表已创建
	if !st.db.Migrator().HasTable(&Translation{}) {
		t.Error("translations 表未创建")
	}
}

func TestPutCache_AndGetCached(t *testing.T) {
	st := newTestStore(t)

	entry := &Translation{
		SourceHash: "abc123",
		TargetLang: "zh-CN",
		SourceText: "# Hello\n\nWorld.",
		Translated: "# 你好\n\n世界。",
		Model:      "gpt-4o",
		TokensUsed: 42,
	}

	// 写入
	if err := st.PutCache(entry); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	// 读取
	cached, err := st.GetCached("abc123", "zh-CN")
	if err != nil {
		t.Fatalf("读取缓存失败: %v", err)
	}
	if cached == nil {
		t.Fatal("缓存未命中")
	}
	if cached.Translated != "# 你好\n\n世界。" {
		t.Errorf("翻译内容不匹配: %q", cached.Translated)
	}
	if cached.TokensUsed != 42 {
		t.Errorf("Tokens: 期望 42, 实际 %d", cached.TokensUsed)
	}
}

func TestGetCached_NotFound(t *testing.T) {
	st := newTestStore(t)

	cached, err := st.GetCached("nonexistent", "ja")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if cached != nil {
		t.Error("未命中应返回 nil")
	}
}

func TestPutCache_Upsert(t *testing.T) {
	st := newTestStore(t)

	// 第一次写入
	st.PutCache(&Translation{
		SourceHash: "hash1",
		TargetLang: "ko",
		SourceText: "Original",
		Translated: "번역 1",
		Model:      "gpt-4o",
		TokensUsed: 10,
	})

	// 第二次写入 (更新)
	st.PutCache(&Translation{
		SourceHash: "hash1",
		TargetLang: "ko",
		SourceText: "Original Updated",
		Translated: "번역 2",
		Model:      "gpt-4o-mini",
		TokensUsed: 20,
	})

	cached, _ := st.GetCached("hash1", "ko")
	if cached.Translated != "번역 2" {
		t.Errorf("Upsert 应该更新: %q", cached.Translated)
	}
}

func TestNewStore_TildeExpansion(t *testing.T) {
	st, err := New("~/test-globalizer-tmp.db")
	if err != nil {
		t.Fatalf("创建 Store 失败: %v", err)
	}
	defer st.Close()

	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, "test-globalizer-tmp.db")

	// 清理测试文件
	defer os.Remove(expectedPath)
	defer os.Remove(expectedPath + "-journal")
	defer os.Remove(expectedPath + "-wal")
	defer os.Remove(expectedPath + "-shm")

	// 验证文件确实在展开后的路径创建
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("数据库文件未在预期路径创建: %s", expectedPath)
	}
}

func TestCache_MultiLangIsolation(t *testing.T) {
	st := newTestStore(t)

	// 写入 zh-CN 缓存
	st.PutCache(&Translation{
		SourceHash: "hash_multi",
		TargetLang: "zh-CN",
		SourceText: "Hello",
		Translated: "你好",
		Model:      "gpt-4o",
	})

	// 写入 ja 缓存 (相同 source_hash, 不同 target_lang)
	st.PutCache(&Translation{
		SourceHash: "hash_multi",
		TargetLang: "ja",
		SourceText: "Hello",
		Translated: "こんにちは",
		Model:      "gpt-4o",
	})

	// 各自独立读取
	zh, _ := st.GetCached("hash_multi", "zh-CN")
	if zh == nil || zh.Translated != "你好" {
		t.Error("zh-CN 缓存读取错误")
	}

	ja, _ := st.GetCached("hash_multi", "ja")
	if ja == nil || ja.Translated != "こんにちは" {
		t.Error("ja 缓存读取错误")
	}
}

func TestStore_DB(t *testing.T) {
	st := newTestStore(t)
	db := st.DB()
	if db == nil {
		t.Fatal("DB() 返回 nil")
	}
	// 验证底层 GORM 实例可用
	var count int64
	db.Model(&Translation{}).Count(&count)
	if count < 0 {
		t.Error("计数不应为负")
	}
}

func TestStore_CloseAndReopen(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	st, err := New(dbPath)
	if err != nil {
		t.Fatalf("创建 Store 失败: %v", err)
	}

	st.PutCache(&Translation{
		SourceHash: "hash_close",
		TargetLang: "zh-CN",
		SourceText: "Hello",
		Translated: "你好",
		Model:      "gpt-4o",
		TokensUsed: 10,
	})

	if err := st.Close(); err != nil {
		t.Fatalf("关闭失败: %v", err)
	}

	// 重新打开同一个数据库文件
	st2, err := New(dbPath)
	if err != nil {
		t.Fatalf("重新打开失败: %v", err)
	}
	defer st2.Close()

	cached, _ := st2.GetCached("hash_close", "zh-CN")
	if cached == nil {
		t.Error("关闭后重新打开，缓存应仍然存在")
	}
	if cached.Translated != "你好" {
		t.Errorf("翻译内容不匹配: %q", cached.Translated)
	}
}
