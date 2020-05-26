package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

func NewMySQL(c Config) (*xorm.Engine, error) {
	dsn := buildDSN(c)
	engine, err := xorm.NewEngine("mysql", dsn)

	if err != nil {
		return nil, err
	}

	engine.TZLocation, err = time.LoadLocation(c.TZLocation)
	if err != nil {
		return nil, err
	}
	engine.DatabaseTZ, err = time.LoadLocation(c.DataBaseTZ)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		engine.ShowSQL(true)
		engine.Logger().SetLevel(log.LOG_DEBUG)
	}

	engine.SetTableMapper(names.GonicMapper{})
	engine.SetColumnMapper(names.GonicMapper{})

	engine.SetMaxIdleConns(c.PoolSize)
	engine.SetMaxOpenConns(c.PoolSize)
	engine.SetConnMaxLifetime(60 * time.Second)

	return engine, nil
}

func buildDSN(c Config) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Name, c.Charset)
}
