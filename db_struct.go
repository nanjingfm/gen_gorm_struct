package structgenerator

import (
	"database/sql"
	"errors"
	"fmt"
	"go/format"
	"strconv"
	"strings"
)

type DbStruct struct {
	db         *sql.DB
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
}

func NewDbStruct(dbName string, dbUser string, dbPassword string, dbHost string, dbPort int) (ds DbStruct, err error) {
	ds.dbName = dbName
	ds.dbUser = dbUser
	ds.dbPassword = dbPassword
	ds.dbHost = dbHost
	ds.dbPort = dbPort
	err = ds.genConn()
	return
}

func (p *DbStruct) openConn() (db *sql.DB, err error) {
	sqlConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", p.dbUser, p.dbPassword, p.dbHost, strconv.Itoa(p.dbPort), p.dbName)
	db, err = sql.Open("mysql", sqlConn)
	if err != nil {
		return
	}

	return
}

func (p *DbStruct) genConn() (err error) {
	if p.db == nil {
		p.db, err = p.openConn()
	}
	return
}

// getAllTableName 获取数据库所有的表
func (p *DbStruct) getAllTableName() (tableNameList []string, err error) {
	var (
		rows *sql.Rows
	)
	tableQuery := fmt.Sprintf("SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s'", p.dbName)
	rows, err = p.db.Query(tableQuery)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return
		}
		tableNameList = append(tableNameList, tableName)
	}
	return
}

// GenDatabase 生成整个库
func (p *DbStruct) GenDatabase() (result string, err error) {
	var (
		tableStruct []string
		ret         string
	)
	tableNameList, err := p.getAllTableName()
	if err != nil {
		return
	}

	for _, tableName := range tableNameList {
		ret, err = p.GenTable(tableName)
		if err != nil {
			return
		}
		tableStruct = append(tableStruct, ret)
	}
	return strings.Join(tableStruct, " "), err
}

// GenTable 生成单个表
func (p *DbStruct) GenTable(tableName string) (result string, err error) {
	columnAttrs, err := p.getTableColumns(tableName)

	if err != nil {
		p.log("Error in selecting column data information from mysql information schema")
		return
	}
	tResult, err := p.genTableStruct(columnAttrs, tableName, tableName)

	if err != nil {
		p.log("Error in creating struct from json: " + err.Error())
		return
	}

	return string(tResult), err
}

func (p *DbStruct) genTableStruct(columnAttrs []Column, tableName string, structName string) ([]byte, error) {
	var dbTypes string
	dbTypes = p.GenerateMysqlTypes(columnAttrs)
	src := fmt.Sprintf(
		"\n//table %s\ntype %s %s}",
		tableName,
		utils.FmtFieldName(utils.StringifyFirstChar(structName)),
		dbTypes,
	)
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

// getTableColumns 获取数据表所有字段
func (p *DbStruct) getTableColumns(dbTable string) (columns []Column, err error) {
	columnDataTypeQuery := "SELECT " +
		"extra, " +
		"column_name, " +
		"data_type, " + // 字段类型 varchar
		"column_type, " + // 字段类型 varchar(32)
		"column_default, " + // 默认值
		"is_nullable, " + // 是否为空
		"column_key, " + // 主键
		"column_comment " + // 备注
		"FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND " +
		"table_name = ?  " +
		"order by `ORDINAL_POSITION` asc"

	rows, err := p.db.Query(columnDataTypeQuery, p.dbName, dbTable)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows == nil {
		return nil, errors.New("no results returned for table")
	}

	for rows.Next() {
		var column Column
		err = rows.Scan(&column.Extra, &column.ColumnName, &column.DataType, &column.ColumnType, &column.ColumnDefault, &column.NullAble, &column.ColumnKey, &column.ColumnComment)
		if err != nil {
			// sql: Scan error on column index 4, name "column_default": unsupported Scan, storing driver.Value type <nil> into type *string
			// TODO
			err = nil
		}

		columns = append(columns, column)
	}

	return columns, err
}

func (p *DbStruct) log(msg string, data ...interface{}) {
	d := make([]interface{}, 0, len(data)+1)
	d = append(d, msg)
	d = append(d, data...)
	fmt.Println(d)
}

func (p *DbStruct) GenerateMysqlTypes(columnAttrs []Column) string {
	structure := "struct {"

	for _, column := range columnAttrs {
		key := column.ColumnName
		var (
			jsonstr   string   //json和form解析
			fields    []string //字段属性
			valueType string   //字段类型
		)

		//是否允许为空
		nullable := false
		if column.NullAble == "YES" {
			nullable = true
		}

		//数据库字段类型转成go数据类型
		valueType = MysqlTypeToGoType(column.DataType, nullable, true)

		//驼峰字段
		fieldName := utils.FmtFieldName(utils.StringifyFirstChar(key))
		jsonstr = fmt.Sprintf("json:\"%s\" form:\"%s\"", key, key)

		if !nullable {
			fields = append(fields, fmt.Sprintf("type:%s NOT NULL", column.ColumnType))
		} else {
			fields = append(fields, fmt.Sprintf("type:%s", column.ColumnType))
		}

		if column.ColumnKey == "PRI" {
			fields = append(fields, "primary_key")
		}
		if column.Extra != "" {
			fields = append(fields, column.Extra+"")
		}

		if column.ColumnDefault != "" {
			fields = append(fields, fmt.Sprintf("DEFAULT:'%s'", column.ColumnDefault))
		}
		fields = append(fields, fmt.Sprintf("COMMENT:'%s'", column.ColumnComment))

		if len(fields) > 0 {
			structure += fmt.Sprintf("\n\t%s %s `%s gorm:\"%s\"`",
				fieldName,
				valueType,
				jsonstr,
				strings.Join(fields, "; "),
			)

		} else {
			structure += fmt.Sprintf("\n%s %s", fieldName, valueType)
		}
	}

	return structure
}
