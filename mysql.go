package gengormstruct

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func NewMysqlDB(config DbConfig) (mdb MysqlDB, err error) {
	mdb.DbConfig = config
	err = mdb.genConn()
	return
}

type MysqlDB struct {
	db *sql.DB
	DbConfig
}

func (p *MysqlDB) openConn() (db *sql.DB, err error) {
	sqlConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", p.DbUser, p.DbPassword, p.DbHost, strconv.Itoa(p.DbPort), p.DbName)
	db, err = sql.Open("mysql", sqlConn)
	if err != nil {
		return
	}

	return
}

func (p *MysqlDB) genConn() (err error) {
	if p.db == nil {
		p.db, err = p.openConn()
	}
	return
}

// GetTableNames returns a slice of table names in the current database
func (p *MysqlDB) GetTableNames() (tables []string) {
	rows, err := p.db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("Could not show tables: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Could not show tables: %s", err)
		}
		tables = append(tables, name)
	}
	return
}

func (p *MysqlDB) GetTable(tableName string) (table *Table) {
	pkColumns := p.GetPkColumns(tableName)
	columns := p.GetColumns(tableName, pkColumns)
	table = new(Table)
	table.Name = CamelCaseString(tableName)
	table.Columns = columns
	table.Pks = pkColumns
	for _, column := range columns {
		if column.Type == "time.Time" {
			table.ImportTimePkg = true
			break
		}
	}
	return table
}

// GetPkColumns 获取主键
func (p *MysqlDB) GetPkColumns(tableName string) (pks []string) {
	rows, err := p.db.Query(
		`SELECT
			c.constraint_type, u.column_name
		FROM
			information_schema.table_constraints c
		INNER JOIN
			information_schema.key_column_usage u ON c.constraint_name = u.constraint_name
		WHERE
			c.table_schema = ? AND u.table_schema = ? AND c.table_name = ? AND u.table_name = ?`,
		p.DbName, p.DbName, tableName, tableName)
	if err != nil {
		log.Fatal("Could not query INFORMATION_SCHEMA for PK information")
	}
	defer rows.Close()
	for rows.Next() {
		var constraintTypeBytes, columnNameBytes []byte
		if err := rows.Scan(&constraintTypeBytes, &columnNameBytes); err != nil {
			log.Fatal("Could not read INFORMATION_SCHEMA for PK information")
		}
		constraintType, columnName := string(constraintTypeBytes), string(columnNameBytes)
		if constraintType == "PRIMARY KEY" {
			pks = append(pks, columnName)
		}
	}
	return pks
}

// GetColumns retrieves columns details from
// information_schema and fill in the Column struct
func (p *MysqlDB) GetColumns(tableName string, pkColumns []string) (columns []*Column) {
	// retrieve columns
	colDefRows, err := p.db.Query(
		`SELECT
			column_name, data_type, column_type, is_nullable, column_default, extra, column_comment 
		FROM
			information_schema.columns
		WHERE
			table_schema = database() AND table_name = ?`,
		tableName)
	if err != nil {
		log.Fatalf("Could not query the database: %s", err)
	}
	defer colDefRows.Close()

	for colDefRows.Next() {
		// datatype as bytes so that SQL <null> values can be retrieved
		var colNameBytes, dataTypeBytes, columnTypeBytes, isNullableBytes, columnDefaultBytes, extraBytes, columnCommentBytes []byte
		if err := colDefRows.Scan(&colNameBytes, &dataTypeBytes, &columnTypeBytes, &isNullableBytes, &columnDefaultBytes, &extraBytes, &columnCommentBytes); err != nil {
			log.Fatal("Could not query INFORMATION_SCHEMA for column information")
		}
		colName, dataType, columnType, isNullable, columnDefault, extra, columnComment :=
			string(colNameBytes), string(dataTypeBytes), string(columnTypeBytes), string(isNullableBytes), string(columnDefaultBytes), string(extraBytes), string(columnCommentBytes)

		// create a column
		col := new(Column)
		col.Name = CamelCaseString(colName)
		col.Type = MysqlTypeToGoType(FilterColumnTypeSize(columnType))

		// Tag info
		tag := new(OrmTag)
		tag.Column = colName
		tag.Comment = columnComment
		tag.Default = columnDefault
		tag.Type = columnType
		tag.HasDefault = columnDefaultBytes != nil
		if StringInSlice(colName, pkColumns) {
			if extra == "auto_increment" {
				tag.AutoPk = true
			} else {
				tag.Pk = true
			}
			tag.HasDefault = false
		} else {
			if isNullable == "YES" {
				tag.Null = true
			}
			if StringInSlice(dataType, []string{"date", "datetime", "timestamp", "time"}) {
				tag.Type = dataType
				//check auto_now, auto_now_add
				if columnDefault == "CURRENT_TIMESTAMP" && extra == "on update CURRENT_TIMESTAMP" {
					tag.AutoNow = true
				} else if columnDefault == "CURRENT_TIMESTAMP" {
					tag.AutoNowAdd = true
				}
				tag.HasDefault = false
			}
		}

		col.Tag = tag
		columns = append(columns, col)
	}
	return
}
