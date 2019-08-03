package source

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/template"

	configStruct "github.com/HencoSmith/graphql-example-go/config/struct"

	"github.com/doug-martin/goqu/v8"
)

func getConnectionString(config configStruct.Configuration) (string, error) {
	// Lookup details defined in environment variables and overwrite config values
	password := config.Database.Password
	envPassword := os.Getenv("POSTGRES_PASSWORD")
	if len(envPassword) != 0 {
		password = envPassword
	}
	user := config.Database.User
	envUser := os.Getenv("POSTGRES_USER")
	if len(envUser) != 0 {
		user = envUser
	}
	DBName := config.Database.Name
	envDBName := os.Getenv("POSTGRES_DB")
	if len(envDBName) != 0 {
		DBName = envDBName
	}

	input := configStruct.DatabaseConfiguration{
		Name:     DBName,
		User:     user,
		Password: password,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		SSL:      config.Database.SSL,
	}

	connectionString := `postgresql://{{.User}}:{{.Password}}@{{.Host}}:{{.Port}}/{{.Name}}?sslmode={{.SSL}}`
	templateString := template.Must(template.New("connectionString").Parse(connectionString))
	var stringParsed bytes.Buffer
	if errExecute := templateString.Execute(&stringParsed, input); errExecute != nil {
		return "", errExecute
	}

	return stringParsed.String(), nil
}

// ConnectToDB attempts to connect to the database and returns a pointer to the database
// along with an error if applicable
func ConnectToDB(config configStruct.Configuration) (*sql.DB, error) {
	connStr, err := getConnectionString(config)
	if err != nil {
		return nil, err
	}

	fmt.Println("connecting to: " + connStr)
	return sql.Open("postgres", connStr)
}

// InitTables - Create and load DB tables with data
// returns an error or nil if no error ocurred
func InitTables(dialect goqu.DialectWrapper, db *sql.DB) error {
	fmt.Println("initializing DB...")
	// Create the Products table if it does not already exist
	createTable, createErr := db.Query(`
	CREATE TABLE IF NOT EXISTS public.products
	(
		id bigserial NOT NULL,
		name character varying(128) NOT NULL,
		info text,
		price numeric NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT products_name_key_unique UNIQUE (name)
	)
	WITH (
		OIDS = FALSE
	);

	ALTER TABLE public.products
		OWNER to "user";
	`)
	if createErr != nil {
		return createErr
	}
	defer createTable.Close()

	// Build Table Seed
	insertDialect := goqu.Insert("products").Rows(
		goqu.Record{
			"name":  "Chicha Morada",
			"info":  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
			"price": 7.99,
		},
		goqu.Record{
			"name":  "Chicha de jora",
			"info":  "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
			"price": 5.95,
		},
		goqu.Record{
			"name":  "Pisco",
			"info":  "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
			"price": 9.95,
		},
	)
	insertQuery, _, toSQLErr := insertDialect.ToSQL()
	if toSQLErr != nil {
		return toSQLErr
	}

	insertRes, insertErr := db.Query(insertQuery)
	// Ignore constraint errors, which means that the data has already been inserted
	if insertErr != nil {
		if !strings.Contains(insertErr.Error(), `unique constraint "products_name_key_unique"`) {
			return insertErr
		}
	} else {
		defer insertRes.Close()
	}

	fmt.Println("OK")
	return nil
}
