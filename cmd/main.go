package cmd

import (
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	structName string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     int
	dbDatabase string
	table      string
	//h          bool
)

func init() {
	//flag.BoolVar(&h, "h", false, "show help")
	flag.StringVar(&structName, "s", "dbstruct", "struct name")
	flag.StringVar(&dbUser, "u", "root", "db user")
	flag.StringVar(&dbPassword, "p", "root@dev", "db password")
	flag.StringVar(&dbHost, "h", "172.172.177.20", "db host")
	flag.IntVar(&dbPort, "P", 33066, "db port")
	flag.StringVar(&dbDatabase, "d", "dynamic", "db name")
	flag.StringVar(&table, "t", "", "table name")
}

func main() {
	var (
		result string
		err    error
		sg     DbStruct
	)

	flag.Parse()

	sg, err = NewDbStruct(dbDatabase, dbUser, dbPassword, dbHost, dbPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	if table != "" {
		result, err = sg.GenTable(table)
	} else {
		result, err = sg.GenDatabase()
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
