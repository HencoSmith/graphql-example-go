package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/doug-martin/goqu/v8"
	_ "github.com/doug-martin/goqu/v8/dialect/postgres"

	"github.com/HencoSmith/graphql-example-go/graphql/movies"
	source "github.com/HencoSmith/graphql-example-go/source"
)

func main() {
	// Read configuration file
	config := source.GetConfig(".")

	// Connect to the database
	db, errConnect := source.ConnectToDB(config)
	if errConnect != nil {
		log.Fatal(errConnect)
	}

	// Lookup the query builder dialect
	dialect := goqu.Dialect("postgres")

	// Setup DB tables and load with data if applicable
	errInit := source.InitTables(dialect, db)
	if errInit != nil {
		log.Fatal(errInit)
	}

	// Bind Queries
	var queryType = movies.Queries(dialect, db)
	// Bind Mutations
	var mutationType = movies.Mutations(dialect, db)

	// Generate the schema
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	// Handle Playground Hosting
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// GraphQL endpoint
	http.Handle("/graphql", h)

	// Prisma GraphQL playground
	http.Handle("/playground/", http.StripPrefix("/playground/", http.FileServer(http.Dir("views"))))

	// Server startup
	fmt.Println("Server is running on port " + config.Server.Port)
	http.ListenAndServe(":"+config.Server.Port, nil)
}
