package mysqlstarterfx

import (
	"fmt"
	"log/slog"

	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	"gorm.io/gorm"
)

type DbLoaderR func(moduleName string, itemID string, item map[string]any) (*gorm.DB, error)
type DbLoader func(moduleName string, itemID string, item map[string]any)

type MysqlModule struct {
	name     string
	load     any
	logLevel string
}

func NewMysqlModule(name string, load any, logLevel string) *MysqlModule {
	return &MysqlModule{name: name, load: load, logLevel: logLevel}
}

func (m MysqlModule) Name() string {
	return m.name
}

func (m MysqlModule) _load(id string, cfg map[string]any) {
	if id == "" {
		slog.Error("attribute id cannot be empty", "module", m.Name())
		panic(fmt.Errorf("[%s] id empty", m.Name()))
	}

	db, err, provided := m.callLoad(id, cfg)
	if err != nil {
		slog.Error("load error", "module", m.Name(), "id", id, "err", err)
		panic(err)
	}

	if provided {
		slog.Info("loading module successfully", "module", m.Name(), "id", id)
		return
	}

	if db == nil {
		slog.Warn("db is nil", "module", m.Name(), "id", id)
		return
	}

	Provide(id, db)
	slog.Info("loading module successfully", "module", m.Name(), "id", id)
}

func (m MysqlModule) Load(value any) {
	if m.logLevel != "" {
		goframeworkgormmysql.SetDbLogLevel(m.logLevel)
	}

	if value == nil {
		return
	}

	items := toStringMap(value)
	if items["id"] != nil {
		id := toString(items["id"])
		m._load(id, items)
		return
	}

	for key, item := range items {
		m._load(key, toStringMap(item))
	}
}

func (m MysqlModule) Close() {
}

func (m MysqlModule) callLoad(id string, cfg map[string]any) (*gorm.DB, error, bool) {
	cfgStr := toStringMapString(cfg)

	switch f := m.load.(type) {
	case DbLoaderR:
		db, err := f(m.Name(), id, cfg)
		return db, err, false
	case func(moduleName string, itemID string, item map[string]any) (*gorm.DB, error):
		db, err := f(m.Name(), id, cfg)
		return db, err, false
	case func(moduleName string, itemID string, item map[string]any) *gorm.DB:
		db := f(m.Name(), id, cfg)
		return db, nil, false
	case func(itemID string, item map[string]any) (*gorm.DB, error):
		db, err := f(id, cfg)
		return db, err, false
	case func(itemID string, item map[string]any) *gorm.DB:
		db := f(id, cfg)
		return db, nil, false
	case func(item map[string]any) (*gorm.DB, error):
		db, err := f(cfg)
		return db, err, false
	case func(item map[string]any) *gorm.DB:
		db := f(cfg)
		return db, nil, false
	case DbLoader:
		f(m.Name(), id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(moduleName string, itemID string, item map[string]any):
		f(m.Name(), id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(itemID string, item map[string]any):
		f(id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(item map[string]any):
		f(cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case nil:
		var err error
		if cfgStr["dsn"] != "" {
			err = goframeworkgormmysql.AddMysqlInstanceWithDsn(id, cfgStr["dsn"])
		} else {
			err = goframeworkgormmysql.AddMysqlInstance(id, cfgStr)
		}
		return nil, err, false
	default:
		return nil, fmt.Errorf("unsupported load callback type: %T", m.load), false
	}
}

func getProvided(id string) *gorm.DB {
	mu.RLock()
	db := dbs[id]
	mu.RUnlock()
	return db
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toStringMap(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	switch m := v.(type) {
	case map[string]any:
		return m
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			out[toString(k)] = val
		}
		return out
	default:
		return map[string]any{}
	}
}

func toStringMapString(v any) map[string]string {
	if v == nil {
		return map[string]string{}
	}
	switch m := v.(type) {
	case map[string]string:
		return m
	case map[string]any:
		out := make(map[string]string, len(m))
		for k, val := range m {
			out[k] = toString(val)
		}
		return out
	case map[any]any:
		out := make(map[string]string, len(m))
		for k, val := range m {
			out[toString(k)] = toString(val)
		}
		return out
	default:
		return map[string]string{}
	}
}
