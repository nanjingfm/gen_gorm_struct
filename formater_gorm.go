package gengormstruct

import (
	"fmt"
	"strings"
)

func NewGormGenerator(c OutputOption) *GormGenerator {
	g := new(GormGenerator)
	g.config = c
	return g
}

type GormGenerator struct {
	config OutputOption
}

func (p *GormGenerator) FormatTable(table *Table) string {
	data := fmt.Sprintf("type %s struct {\n", table.Name)

	// fill column
	for _, column := range table.Columns {
		data += p.formatColumn(column) + "\n"
	}

	data += "}\n"
	return data
}

func (p *GormGenerator) formatColumn(column *Column) string {
	tag := column.Tag
	var (
		tagList []string
		tagStr  string
	)
	if tag.AutoPk {
		tagList = append(tagList, fmt.Sprintf("type:%s auto_increment", tag.Type))
	} else {
		tagList = append(tagList, "type:"+tag.Type)
	}
	if !tag.Null {
		tagList = append(tagList, "NOT NULL")
	}
	if tag.HasDefault {
		tagList = append(tagList, fmt.Sprintf("DEFAULT:'%s'", tag.Default))
	}
	if tag.AutoNow {
		tagList = append(tagList, "DEFAULT: CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	}
	if tag.AutoNowAdd {
		tagList = append(tagList, "DEFAULT: CURRENT_TIMESTAMP")
	}
	if tag.AutoPk || tag.Pk {
		tagList = append(tagList, "primary_key")
	}
	if tag.Comment != "" {
		tagList = append(tagList, fmt.Sprintf("COMMENT:'%s'", tag.Comment))
	}
	if p.config.JsonTag {
		tagStr += fmt.Sprintf("json:\"%s\" ", tag.Column)
	}
	if p.config.FormTag {
		tagStr += fmt.Sprintf("form:\"%s\" ", tag.Column)
	}
	tagStr += fmt.Sprintf("gorm:\"%s\"", strings.Join(tagList, "; "))
	return fmt.Sprintf("%s %s `%s`", column.Name, column.Type, tagStr)
}
