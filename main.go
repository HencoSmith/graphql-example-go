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

	"github.com/HencoSmith/graphql-example-go/graphql/products"
	"github.com/HencoSmith/graphql-example-go/models"
	source "github.com/HencoSmith/graphql-example-go/source"
)

// Refer to: https://github.com/graphql-go/graphql/tree/master/examples/crud
// Reworking into a project template

var productArr []models.Product

func initProductsData(p *[]models.Product) {
	product1 := models.Product{ID: 1, Name: "Chicha Morada", Info: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)", Price: 7.99}
	product2 := models.Product{ID: 2, Name: "Chicha de jora", Info: "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)", Price: 5.95}
	product3 := models.Product{ID: 3, Name: "Pisco", Info: "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)", Price: 9.95}
	*p = append(*p, product1, product2, product3)
}

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

	// Primary data initialization
	initProductsData(&productArr)

	// Bind Product Queries
	var queryType = products.Queries(&productArr, dialect, db)
	// Bind Product Mutations
	var mutationType = products.Mutations(&productArr)

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
