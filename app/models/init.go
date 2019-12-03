package models

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"

	// Required to interact with postgres dbs
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "john"
	password = "password"
	dbname   = "gameshelf"
)

var db *sql.DB
var dbmap *gorp.DbMap

// InitDB initializes the database
func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{LowercaseFields: false}}

	if err != nil {
		panic(err)
	}

	// defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	dbmap.AddTableWithName(Game{}, "games").SetKeys(true, "ID")
	dbmap.AddTableWithName(User{}, "users").SetKeys(false, "Username")

}
