package gengormstruct

import (
	"fmt"
)

// GenGormTable 生成单表
func GenGormTable(config DbConfig, option OutputOption, tableName string) (err error)  {
	db, err := NewMysqlDB(config)
	if err != nil {
		return
	}
	table := db.GetTable(tableName)
	code := NewGormGenerator(option).FormatTable(table)
	_, _ = fmt.Fprint(option.GetOutputMode(), FormatSourceCode(code))
	return
}

// GenOrmDatabase 生成整个数据库
func GenOrmDatabase(config DbConfig, option OutputOption) (err error)  {
	db, err := NewMysqlDB(config)
	if err != nil {
		return
	}
	tableList := db.GetTableNames()
	for _, tableName := range tableList {
		err = GenGormTable(config, option, tableName)
		if err != nil {
			return
		}
	}
	return nil
}
