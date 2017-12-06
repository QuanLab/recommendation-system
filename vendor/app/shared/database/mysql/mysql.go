package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

var (
	DB *sql.DB
	DB2 *sql.DB
)

//datasource name
func DSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

//open connection to MySQL database, it 's self contains pool init
func Connect(d MySQLInfo) {
	var err error
	DB, err = sql.Open("mysql", DSN(d))
	if err != nil {
		log.Printf("Cannot connect to MySQL server %s", err)
	}
}

func ConnectDB2(d MySQLInfo) {
	var err error
	DB2, err = sql.Open("mysql", DSN(d))
	if err != nil {
		log.Printf("Cannot connect to MySQL server %s", err)
	}
}
