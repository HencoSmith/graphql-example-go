package source

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
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
	// Create the Products table if it does not already exist
	createRows, createErr := db.Query(`
	CREATE TABLE IF NOT EXISTS public.products
	(
		id bigserial NOT NULL,
		name character varying(128) NOT NULL,
		info text,
		price numeric NOT NULL,
		PRIMARY KEY (id)
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
	defer createRows.Close()

	// TODO: complete insert of data
	// dialectString := dialect.From("products").Insert(
	// 	goqu.Record{
	// 		"name":  "Chicha Morada",
	// 		"info":  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
	// 		"price": 7.99,
	// 	},
	// 	goqu.Record{
	// 		"name":  "Chicha de jora",
	// 		"info":  "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
	// 		"price": 5.95,
	// 	},
	// 	goqu.Record{
	// 		"name":  "Pisco",
	// 		"info":  "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
	// 		"price": 9.95,
	// 	},
	// )
	// query, args, insertErr := dialectString.toSQL()
	// if insertErr != nil {
	// 	return insertErr
	// }

	// fmt.Println(query, args)

	// var names []string
	// for createRows.Next() {
	// 	var name string
	// 	if err := createRows.Scan(&name); err != nil {
	// 		return err
	// 	}
	// 	names = append(names, name)
	// }

	// fmt.Println(names)
	return nil
}
