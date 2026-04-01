package mysql_starter

import (
	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	"github.com/spf13/cast"
	"log/slog"
	"os"
)

func HasMysqlInstance(db string) bool {
	return goframeworkgormmysql.HasMysqlInstance(db)
}

func CloseMysqlInstance(db string) {
	goframeworkgormmysql.RemoveMysqlInstance(db)
}

type MysqlModule struct {
	name     string
	load     func(moduleName string, itemId string, item map[string]string)
	logLevel string
}

func NewMysqlModule(name string, load func(moduleName string, itemId string, item map[string]string), logLevel string) *MysqlModule {
	return &MysqlModule{name, load, logLevel}
}

func (m MysqlModule) Name() string {
	return m.name
}

func (m MysqlModule) _load(id string, cfg map[string]string) {
	if id == "" {
		slog.Error("attribute id cannot be empty", "module", m.Name())
		os.Exit(1)
		return
	}

	if cfg["dsn"] != "" {
		if err := goframeworkgormmysql.AddMysqlInstanceWithDsn(id, cfg["dsn"]); err != nil {
			slog.Error("initializing mysql with dsn failed", "module", m.Name(), "id", id, "err", err)
			os.Exit(1)
			return
		}
	} else {
		if err := goframeworkgormmysql.AddMysqlInstance(id, cfg); err != nil {
			slog.Error("initializing mysql failed", "module", m.Name(), "id", id, "err", err)
			os.Exit(1)
			return
		}
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		slog.Debug("triggering custom loader completion", "module", m.Name())
	}

	slog.Info("loading module successfully", "module", m.Name(), "id", id)
}

func (m MysqlModule) Load(value interface{}) {

	if m.logLevel != "" {
		goframeworkgormmysql.SetDbLogLevel(m.logLevel)
	}

	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}
}

func (m MysqlModule) Close() {
}
