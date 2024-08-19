package mysql_starter

import (
	"github.com/kordar/gocfg"
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
	GroupName string
}

func (m MysqlModule) Name() string {
	return "mysql_starter"
}

func (m MysqlModule) Load(value interface{}) {

	logLevel := gocfg.GetSystemValue("gorm_log_level", m.GroupName)
	if logLevel != "" {
		goframeworkgormmysql.SetDbLogLevel(logLevel)
	}

	items := cast.ToStringMap(value)
	for key, item := range items {
		section := cast.ToStringMapString(item)
		if section["dsn"] != "" {
			if err := goframeworkgormmysql.AddMysqlInstanceWithDsn(key, section["dsn"]); err != nil {
				logger.Fatalf("[mysql_starter] 初始化mysql异常，err=%v", err)
			}
		} else {
			if err := goframeworkgormmysql.AddMysqlInstance(key, section); err != nil {
				logger.Fatalf("[mysql_starter] 初始化mysql异常，err=%v", err)
			}
		}
	}
}

func (m MysqlModule) Close() {
}
