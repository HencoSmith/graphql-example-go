package products

import (
	"fmt"

	"database/sql"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
)

// Queries - all GraphQL queries related to products
func Queries(products *[]models.Product, dialect goqu.DialectWrapper, db *sql.DB) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				/* Get (read single product by id)
				http://localhost:8080/graphql?query={product(id:1){name,info,price}}
				*/
				"product": &graphql.Field{
					Type:        ProductType,
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
							for _, product := range *products {
								if int(product.ID) == id {
									return product, nil
								}
							}
						}
						return nil, nil
					},
				},
				/* Get (read) product list
				http://localhost:8080/graphql?query={list{id,name,info,price}}
				*/
				"list": &graphql.Field{
					Type:        graphql.NewList(ProductType),
					Description: "Get product list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						// TODO: finish query builder
						dialectString := dialect.From("products").Where(goqu.Ex{"id": 10})
						query, args, err := dialectString.ToSQL()
						if err != nil {
							fmt.Println("Failed to generate query string", err.Error())
						} else {
							fmt.Println(query, args)
						}

						rows, err := db.Query(query)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println(rows)
						}
						// END

						return *products, nil
					},
				},
			},
		},
	)
}
