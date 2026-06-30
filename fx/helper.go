package mysqlstarterfx

import (
	"fmt"
	"log/slog"
	"sync"

	"gorm.io/gorm"
)

var (
	dbs = make(map[string]*gorm.DB)
	mu  sync.RWMutex
)

func Get(name string) *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	db, ok := dbs[name]
	if !ok {
		slog.Error("mysql db not exist", "name", name)
		panic(fmt.Errorf("mysql db %s not exist", name))
	}
	return db
}

func Provide(name string, db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	dbs[name] = db
}
