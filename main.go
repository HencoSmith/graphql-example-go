package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/mnmtanish/go-graphiql"

	"github.com/HencoSmith/graphql-example-go/graphql/products"
	"github.com/HencoSmith/graphql-example-go/models"
)

// Refer to: https://github.com/graphql-go/graphql/tree/master/examples/crud
// Reworking into a project template

var productArr []models.Product

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("\nerrors: %v", result.Errors)
	}
	return result
}

func initProductsData(p *[]models.Product) {
	product1 := models.Product{ID: 1, Name: "Chicha Morada", Info: "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)", Price: 7.99}
	product2 := models.Product{ID: 2, Name: "Chicha de jora", Info: "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)", Price: 5.95}
	product3 := models.Product{ID: 3, Name: "Pisco", Info: "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)", Price: 9.95}
	*p = append(*p, product1, product2, product3)
}

func main() {
	// Primary data initialization
	initProductsData(&productArr)

	// Bind Product Queries
	var queryType = products.Queries(&productArr)
	// Bind Product Mutations
	var mutationType = products.Mutations(&productArr)

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	// GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Graphiql Does not post any variables only an query string e.g.
		// r.URL.Query().Get("query") = "" when graphiql runs
		fmt.Printf("%v %v %v", r.Method, r.URL, r.Proto)
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	// GraphQL Playground
	http.HandleFunc("/graphiql", graphiql.ServeGraphiQL)

	// Prisma GraphQL playground
	http.Handle("/playground/", http.StripPrefix("/playground/", http.FileServer(http.Dir("views"))))

	// Server startup
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
