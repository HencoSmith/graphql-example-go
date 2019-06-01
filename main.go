package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/HencoSmith/graphql-example-go/graphql/products"
	"github.com/HencoSmith/graphql-example-go/models"
)

// Refer to: https://github.com/graphql-go/graphql/tree/master/examples/crud
// Reworking into a project template

var productArr []models.Product

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read single product by id)
			http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			"product": &graphql.Field{
				Type:        products.ProductType,
				Description: "Get product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find product
						for _, product := range productArr {
							if int(product.ID) == id {
								return product, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(products.ProductType),
				Description: "Get product list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return productArr, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		"create": &graphql.Field{
			Type:        products.ProductType,
			Description: "Create new product",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"info": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				product := models.Product{
					ID:    int64(rand.Intn(100000)), // generate random ID
					Name:  params.Args["name"].(string),
					Info:  params.Args["info"].(string),
					Price: params.Args["price"].(float64),
				}
				productArr = append(productArr, product)
				return product, nil
			},
		},

		/* Update product by id
		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
		*/
		"update": &graphql.Field{
			Type:        products.ProductType,
			Description: "Update product by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"info": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				name, nameOK := params.Args["name"].(string)
				info, infoOK := params.Args["info"].(string)
				price, priceOK := params.Args["price"].(float64)
				product := models.Product{}
				for i, p := range productArr {
					if int64(id) == p.ID {
						if nameOK {
							productArr[i].Name = name
						}
						if infoOK {
							productArr[i].Info = info
						}
						if priceOK {
							productArr[i].Price = price
						}
						product = productArr[i]
						break
					}
				}
				return product, nil
			},
		},

		/* Delete product by id
		http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
		*/
		"delete": &graphql.Field{
			Type:        products.ProductType,
			Description: "Delete product by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				product := models.Product{}
				for i, p := range productArr {
					if int64(id) == p.ID {
						product = productArr[i]
						// Remove from product list
						productArr = append(productArr[:i], productArr[i+1:]...)
					}
				}

				return product, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
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

	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
