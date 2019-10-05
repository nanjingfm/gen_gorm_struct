package structgenerator

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func getTestInc() (sg DbStruct, err error) {
	var (
		dbUser     string = "root"
		dbPassword string = "root@dev"
		dbHost     string = "10.254.0.88"
		dbPort     int    = 3306
		dbName     string = "marketing_tools"
	)
	sg, err = NewDbStruct(dbName, dbUser, dbPassword, dbHost, dbPort)
	return
}

func TestDbStruct_GenTable(t *testing.T) {
	sg, err := getTestInc()
	assert.Nil(t, err, "getTestInc err")
	result, err := sg.GenTable("activity")
	assert.Nil(t, err, "GenTable err")
	spew.Println(result)
}

func TestDbStruct_GenDatabase(t *testing.T) {
	sg, err := getTestInc()
	assert.Nil(t, err, "getTestInc err")
	result, err := sg.GenDatabase()
	assert.Nil(t, err, "GenDatabase err")
	spew.Println(result)
}
