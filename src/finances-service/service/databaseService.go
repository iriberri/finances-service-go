package service

import (
    "database/sql"
    "fmt"
    "log"
    "strings"

    _ "github.com/lib/pq"
    "github.com/olebedev/config"
)

type DatabaseService interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}

type databaseService struct {
    db *sql.DB
}

var _ DatabaseService = &databaseService{}

// Initialize the `DatabaseService`. Will panic if the connection to the database
// cannot be established.
func NewDatabaseService(conf *config.Config) DatabaseService {
    dbConf, err := conf.Get("database")
    if err != nil {
        panic (err)
    }
    username := dbConf.UString("username")
    password := dbConf.UString("password")
    hostname := dbConf.UString("hostname")
    port := dbConf.UInt("port")
    database := dbConf.UString("database")
    options := dbConf.UMap("options")

    optionsString := toUrlQuery(options)
    connStr := fmt.Sprintf(
        "postgresql://%s:%s@%s:%d/%s%s",
        username, password, hostname, port, database, optionsString)
    log.Printf("Connecting to database using connection string: %s", connStr)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
    return &databaseService{
        db: db,
    }
}

func (db *databaseService) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return db.db.Query(query, args...)
}

func (db *databaseService) QueryRow(query string, args ...interface{}) *sql.Row {
    return db.db.QueryRow(query, args...)
}

func toUrlQuery(m map[string]interface{}) string {
    if len(m) > 0 {
        builder := strings.Builder{}
        builder.WriteRune('?')
        separate := false
        for k, v := range m {
            if separate {
                builder.WriteRune('&')
            }
            builder.WriteString(k)
            builder.WriteRune('=')
            builder.WriteString(fmt.Sprint(v))
            separate = true
        }
        return builder.String()
    }
    return ""
}
