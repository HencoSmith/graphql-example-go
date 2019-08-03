package products

import (
	"database/sql"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
)

// Queries - all GraphQL queries related to products
func Queries(dialect goqu.DialectWrapper, db *sql.DB) *graphql.Object {
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
							dialectString := dialect.From("products").Where(goqu.Ex{
								"id": id,
							})
							query, _, dialectErr := dialectString.ToSQL()
							if dialectErr != nil {
								return nil, dialectErr
							}

							rows, queryErr := db.Query(query)
							if queryErr != nil {
								return nil, queryErr
							}
							defer rows.Close()

							var productsArr []models.Product
							for rows.Next() {
								var productRow = models.Product{}
								scanErr := rows.Scan(&productRow.ID, &productRow.Name, &productRow.Info, &productRow.Price)
								if scanErr != nil {
									return nil, scanErr
								}
								productsArr = append(productsArr, productRow)
							}
							if errRows := rows.Err(); errRows != nil {
								return nil, errRows
							}

							return &productsArr[0], nil
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
						dialectString := dialect.From("products").Order(goqu.C("id").Asc())
						query, _, dialectErr := dialectString.ToSQL()
						if dialectErr != nil {
							return nil, dialectErr
						}

						rows, queryErr := db.Query(query)
						if queryErr != nil {
							return nil, queryErr
						}
						defer rows.Close()

						var productsArr []models.Product
						for rows.Next() {
							var productRow = models.Product{}
							scanErr := rows.Scan(&productRow.ID, &productRow.Name, &productRow.Info, &productRow.Price)
							if scanErr != nil {
								return nil, scanErr
							}
							productsArr = append(productsArr, productRow)
						}
						if errRows := rows.Err(); errRows != nil {
							return nil, errRows
						}

						return &productsArr, nil
					},
				},
			},
		},
	)
}
