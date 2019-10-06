package main

import (
	"flag"
	"fmt"
	gengormstruct "github.com/nanjingfm/gen_gorm_struct"

	_ "github.com/go-sql-driver/mysql"
)

var (
	tableName string
	config    gengormstruct.DbConfig
	option    gengormstruct.OutputOption
)

func init() {
	flag.StringVar(&config.DbUser, "u", "root", "db user")
	flag.StringVar(&config.DbPassword, "p", "123456", "db password")
	flag.StringVar(&config.DbHost, "h", "127.0.0.1", "db host")
	flag.IntVar(&config.DbPort, "P", 3306, "db port")
	flag.StringVar(&config.DbName, "d", "test", "db name")
	flag.StringVar(&tableName, "t", "", "table name")
	flag.BoolVar(&option.JsonTag, "json", false, "with json tag")
	flag.BoolVar(&option.FormTag, "form", false, "with form tag")
}

func main() {
	var err error

	flag.Parse()
	if tableName != "" {
		err = gengormstruct.GenGormTable(config, option, tableName)
	} else {
		err = gengormstruct.GenOrmDatabase(config, option)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

}
