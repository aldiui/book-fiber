package connection

import (
	"book-fiber/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func GetDatabase(conf config.Database) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.Name, conf.Tz)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database ", err)
	}
	err = db.Ping()

	if err != nil {
		log.Fatal("Error connecting to database ", err)
	}

	return db
}
