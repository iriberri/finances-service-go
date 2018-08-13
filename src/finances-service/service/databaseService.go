package service

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "path"
    "regexp"
    "strconv"
    "strings"

    "github.com/GuiaBolso/darwin"
    "github.com/adeynack/finances-service-go/src/finances-service/util"

    // Import the database drive within the service itself to avoid managing DB related
    // concerns in the main package.
    _ "github.com/lib/pq"
)

// DatabaseService allows interactions with the database to be done through an interface,
// so it can easily be mocked for test purposes.
type DatabaseService interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}

type databaseService struct {
    db *sql.DB
}

var _ DatabaseService = &databaseService{}

// NewDatabaseService initializes the default production `DatabaseService`.
// It will panic if the connection to the database cannot be established.
func NewDatabaseService(conf *util.ConfigReader) DatabaseService {
    dbConf := conf.MustGet("database")
    username := dbConf.MustString("username")
    password := dbConf.MustString("password")
    hostname := dbConf.MustString("hostname")
    port := dbConf.MustInt("port")
    database := dbConf.MustString("database")
    schema := dbConf.MustString("schema")
    options := dbConf.MustMap("options")
    options["search_path"] = schema
    optionsString := toURLQuery(options)
    connStr := fmt.Sprintf(
        "postgresql://%s:%s@%s:%d/%s%s",
        username, password, hostname, port, database, optionsString)
    log.Printf("Connecting to database using connection string: %s", connStr)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
        panic(err)
    }

    err = evolveDatabase(db, dbConf.MustGet("evolution"), schema)
    if err != nil {
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

func toURLQuery(m map[string]interface{}) string {
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

func evolveDatabase(db *sql.DB, evConf *util.ConfigReader, schema string) error {

    if err := recreateSchema(db, evConf, schema); err != nil {
        return err
    }

    runAtStartup := evConf.MustBool("run_at_startup")
    if !runAtStartup {
        log.Println("Database evolution are configured not to run.")
        return nil
    }

    steps, err := extractEvolutionFromFiles(evConf)
    if err != nil {
        return err
    }

    infoChan := make(chan darwin.MigrationInfo)
    defer close(infoChan)
    doneChan := make(chan bool)
    defer close(doneChan)
    go logMigrationInfo(infoChan, doneChan)

    driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
    d := darwin.New(driver, steps, infoChan)
    err = d.Migrate()
    infoChan <- darwin.MigrationInfo{Status: -1} // todo: Suggest PR to 'darwin' for a 'done' migration info status
    <-doneChan
    if err != nil {
        return fmt.Errorf("migrating database: %s", err)
    }
    return nil
}

func recreateSchema(db *sql.DB, evConf *util.ConfigReader, schemaName string) error {
    if !evConf.UBool("recreate_schema", false) {
        return nil
    }
    log.Printf("Forcing re-creation of the database schema. HINT: THIS SHOULD NEVER HAPPEN IN PRODUCTION\n")

    // CHECK IF SCHEMA EXISTS
    queryCheckIfSchemaExists := fmt.Sprintf("SELECT schema_name FROM information_schema.schemata WHERE schema_name = '%s'", schemaName)
    res, err := db.Query(queryCheckIfSchemaExists)
    defer res.Close()
    if err != nil {
        return fmt.Errorf("recreate schema: checking if schema exists: %s", err)
    }
    if res.Next() {
        // DROPPING EXISTING SCHEMA
        log.Printf("recreate schema: schema %q exists: dropping\n", schemaName)
        queryDropSchema := fmt.Sprintf("drop schema %s cascade", schemaName)
        _, err := db.Exec(queryDropSchema)
        if err != nil {
            return fmt.Errorf("recreate schema: droping schema %q: %s", schemaName, err)
        }
        log.Printf("recreate schema: schema dropped\n")
    }

    // CREATE SCHEMA
    log.Printf("recreate schema: creating schema %q\n", schemaName)
    queryCreateSchema := fmt.Sprintf("create schema %s", schemaName)
    _, err = db.Exec(queryCreateSchema)
    if err != nil {
        return fmt.Errorf("recreate schema: creating schema %q: %s", schemaName, err)
    }
    return nil
}

func extractEvolutionFromFiles(evConf *util.ConfigReader) ([]darwin.Migration, error) {
    errorSuffix := "evolution folder is expected to contain only files of format \"0001-Name.sql\" where 0001 is an sequential identifier and \"name\" is the description of the evolution step)"
    fileNameRegEx := regexp.MustCompile(`([0-9]*)-([^-]*)(-DEV)?\.sql`)
    steps := make([]darwin.Migration, 0)
    folder := evConf.MustString("scripts_folders")
    includeDev := evConf.UBool("include_dev_scripts", false)
    log.Printf("Scanning folder %q for database evolution scripts\n", folder)
    files, err := ioutil.ReadDir(folder)
    if err != nil {
        return nil, fmt.Errorf("reading folder %q: %s", folder, err)
    }
    for _, file := range files {
        if file.IsDir() {
            return nil, fmt.Errorf("encountered folder %q. %s", file.Name(), errorSuffix)
        }
        matches := fileNameRegEx.FindStringSubmatch(file.Name())
        if matches == nil || len(matches) != 4 {
            return nil, fmt.Errorf("file %q does not match expected format. %s", file.Name(), errorSuffix)
        }
        if !includeDev && len(matches[3]) > 0 {
            log.Printf("%q skipping DEV script file", file.Name())
            continue
        }
        log.Printf("%q reading file", file.Name())
        evolutionNumber, _ := strconv.Atoi(matches[1])
        if evolutionNumber < 1 {
            return nil, fmt.Errorf("file %q: number part should be a positive number. %s", file.Name(), errorSuffix)
        }
        evolutionName := matches[2]
        filePath := path.Join(folder, file.Name())
        evolutionScript, err := ioutil.ReadFile(filePath)
        if err != nil {
            return nil, fmt.Errorf("reading file %s: %s", filePath, err)
        }
        steps = append(steps, darwin.Migration{
            Version:     float64(evolutionNumber),
            Description: evolutionName,
            Script:      string(evolutionScript),
        })
    }
    return steps, nil
}

func logMigrationInfo(infoChan chan darwin.MigrationInfo, doneChan chan bool) {
    count := 0

    for info := range infoChan {
        if info.Status < 0 {
            break
        }
        if count == 0 {
            log.Print("Migration is starting\n")
        }
        count++
        var errorPart string
        if info.Error == nil {
            errorPart = ""
        } else {
            errorPart = fmt.Sprintf(": %s", info.Error)
        }
        log.Printf(
            "Migration %v %q / status: %s%s\n",
            info.Migration.Version,
            info.Migration.Description,
            info.Status,
            errorPart)
    }
    if count == 0 {
        log.Print("No migration to be performed (database is up to date)\n")
    } else {
        log.Printf("Migration is done (%d steps performed)\n", count)
    }
    doneChan <- true
}
