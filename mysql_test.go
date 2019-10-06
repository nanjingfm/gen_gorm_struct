package gengormstruct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getMysqlTestInc() (sg MysqlDB, err error) {
	config := DbConfig{
		DbName:     "test",
		DbUser:     "root",
		DbPassword: "123456",
		DbHost:     "127.0.0.1",
		DbPort:     3306,
	}
	sg, err = NewMysqlDB(config)
	return
}

func TestMysqlDB_GetTableNames(t *testing.T) {
	m, err := getMysqlTestInc()
	assert.Nil(t, err, "getMysqlTestInc err")
	tables := m.GetTableNames()
	assert.NotEqual(t, 0, len(tables))
}

func TestMysqlDB_GetPkColumns(t *testing.T) {
	m, err := getMysqlTestInc()
	assert.Nil(t, err, "getMysqlTestInc err")
	pkColumns := m.GetPkColumns("activity")
	assert.NotEqual(t, 0, len(pkColumns))
}

func TestMysqlDB_GetColumns(t *testing.T) {
	m, err := getMysqlTestInc()
	assert.Nil(t, err, "getMysqlTestInc err")
	columns := m.GetColumns("activity", []string{"id"})
	assert.NotEqual(t, 0, len(columns))
}

func TestMysqlDB_GetTable(t *testing.T) {
	m, err := getMysqlTestInc()
	assert.Nil(t, err, "getMysqlTestInc err")
	table := m.GetTable("activity")
	assert.NotNil(t, table)
	//fmt.Println(FormatSourceCode(NewGormGenerator(OutputOption{}).FormatTable(table)))
}
