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

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	dbmap.AddTableWithName(Game{}, "games").SetKeys(true, "ID")
	dbmap.AddTableWithName(User{}, "users").SetKeys(false, "Username")
	dbmap.AddTableWithName(Match{}, "matches").SetKeys(true, "ID")
	dbmap.AddTableWithName(MatchScore{}, "match_scores").SetKeys(true, "ID")
	dbmap.AddTableWithName(Friend{}, "friends").SetKeys(true, "ID")
	dbmap.AddTableWithName(Group{}, "groups").SetKeys(true, "ID")
	dbmap.AddTableWithName(GroupMember{}, "group_members").SetKeys(true, "ID")
	dbmap.AddTableWithName(GroupGame{}, "group_games").SetKeys(true, "ID")

	dbmap.ExpandSliceArgs = true
}

func findQstring(table string) string {
	return fmt.Sprintf("select * from %s where \"ID\"=$1", table)
}
