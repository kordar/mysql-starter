package mysqlstarterfx

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestMysqlModule_Load_Multi(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test_%d_", time.Now().UnixNano())
	created := map[string]*gorm.DB{}

	m := NewMysqlModule("gorm", func(moduleName string, itemID string, item map[string]any) (*gorm.DB, error) {
		_ = moduleName
		_ = item
		db := &gorm.DB{}
		created[itemID] = db
		return db, nil
	}, "")

	m.Load(map[string]any{
		prefix + "a": map[string]any{"dsn": "a"},
		prefix + "b": map[string]any{"dsn": "b"},
	})

	if Get(prefix+"a") != created[prefix+"a"] {
		t.Fatalf("db mismatch: a")
	}
	if Get(prefix+"b") != created[prefix+"b"] {
		t.Fatalf("db mismatch: b")
	}
}

func TestMysqlModule_Load_SingleByID(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_solo", time.Now().UnixNano())
	var got *gorm.DB

	m := NewMysqlModule("gorm", func(moduleName string, itemID string, item map[string]any) (*gorm.DB, error) {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		got = &gorm.DB{}
		return got, nil
	}, "")

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != got {
		t.Fatalf("db mismatch: solo")
	}
}

func TestMysqlModule_Load_VoidLoader_Provides(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_void", time.Now().UnixNano())
	want := &gorm.DB{}

	m := NewMysqlModule("gorm", func(moduleName string, itemID string, item map[string]any) {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		Provide(itemID, want)
	}, "")

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != want {
		t.Fatalf("db mismatch: void")
	}
}

func TestMysqlModule_Load_NoErrorReturnLoader(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_noerr", time.Now().UnixNano())
	want := &gorm.DB{}

	m := NewMysqlModule("gorm", func(moduleName string, itemID string, item map[string]any) *gorm.DB {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		return want
	}, "")

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != want {
		t.Fatalf("db mismatch: noerr")
	}
}
