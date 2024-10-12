package mysql_starter

import (
	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	logger "github.com/kordar/gologger"
	"github.com/spf13/cast"
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
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	if cfg["dsn"] != "" {
		if err := goframeworkgormmysql.AddMysqlInstanceWithDsn(id, cfg["dsn"]); err != nil {
			logger.Fatalf("[%s] initializing mysql: %v", m.Name(), err)
			return
		}
	} else {
		if err := goframeworkgormmysql.AddMysqlInstance(id, cfg); err != nil {
			logger.Fatalf("[%s] initializing mysql: %v", m.Name(), err)
			return
		}
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
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
