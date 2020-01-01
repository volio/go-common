package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

func NewMysql(c Config) (*xorm.Engine, error) {
	dsn := buildDSN(c)
	engine, err := xorm.NewEngine("mysql", dsn)

	if err != nil {
		return nil, err
	}

	engine.TZLocation, _ = time.LoadLocation(c.Location)

	if c.Debug {
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
		engine.Logger().SetLevel(core.LOG_DEBUG)
	}

	engine.SetTableMapper(core.GonicMapper{})
	engine.SetColumnMapper(core.GonicMapper{})

	engine.SetMaxIdleConns(c.PoolSize)
	engine.SetMaxOpenConns(c.PoolSize)
	engine.SetConnMaxLifetime(60 * time.Second)

	return engine, nil
}

func buildDSN(c Config) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s", c.Username, c.Password, c.Host, c.Port, c.Name, c.Charset)
}
