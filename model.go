package gengormstruct

import (
	"io"
	"os"
)

var defaultOutputMode = os.Stdout

type Table struct {
	Name          string
	Pks           []string
	Columns       []*Column
	ImportTimePkg bool
}

type Column struct {
	Name string
	Type string
	Tag  *OrmTag
}

type OrmTag struct {
	AutoPk     bool   // 是否是自增主键
	Pk         bool   // 是否是主键
	Null       bool   // 是否允许null
	Column     string // mysql字段名
	AutoNow    bool   // updated_at
	AutoNowAdd bool   // created_at
	Type       string // mysql 字段类型
	Default    string
	HasDefault bool // 是否you默认值
	Comment    string
}

type DbConfig struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     int
}

type OutputOption struct {
	JsonTag bool
	FormTag bool
	OutputMode string
}

func (p OutputOption) GetOutputMode() io.WriteCloser {
	if p.OutputMode == "" {
		return defaultOutputMode
	}

	return defaultOutputMode
}