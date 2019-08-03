package products

import (
	"database/sql"
	"errors"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
)

func findProduct(dialect goqu.DialectWrapper, db *sql.DB, expression goqu.Ex) (*models.Product, error) {
	// Find product inserted (ID is generated by DB hence lookup is required)
	dialectString := dialect.From("products").Where(expression)
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

// Mutations - all GraphQL mutations related to products
func Mutations(dialect goqu.DialectWrapper, db *sql.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* Create new product item
			http://localhost:8080/graphql?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
			*/
			"create": &graphql.Field{
				Type:        ProductType,
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
					// Insert the new product
					insertDialect := goqu.Insert("products").Rows(
						goqu.Record{
							"name":  params.Args["name"].(string),
							"info":  params.Args["info"].(string),
							"price": params.Args["price"].(float64),
						},
					)
					insertQuery, _, toSQLErr := insertDialect.ToSQL()
					if toSQLErr != nil {
						return nil, toSQLErr
					}

					insertRes, insertErr := db.Query(insertQuery)
					if insertErr != nil {
						return nil, insertErr
					}

					defer insertRes.Close()

					product, findErr := findProduct(dialect, db, goqu.Ex{
						// Name has a unique constraint
						"name": params.Args["name"].(string),
					})

					return product, findErr
				},
			},

			/* Update product by id
			   http://localhost:8080/graphql?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
			*/
			"update": &graphql.Field{
				Type:        ProductType,
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

					if !nameOK {
						return nil, errors.New("Invalid name")
					}
					if !infoOK {
						return nil, errors.New("Invalid info")
					}
					if !priceOK {
						return nil, errors.New("Invalid price")
					}

					// Update the existing product
					updateDialect := goqu.Update("products").Set(
						goqu.Record{
							"name":  name,
							"info":  info,
							"price": price,
						},
					).Where(goqu.Ex{
						"id": id,
					})
					updateQuery, _, toSQLErr := updateDialect.ToSQL()
					if toSQLErr != nil {
						return nil, toSQLErr
					}

					updateRes, updateErr := db.Query(updateQuery)
					if updateErr != nil {
						return nil, updateErr
					}

					defer updateRes.Close()

					product, findErr := findProduct(dialect, db, goqu.Ex{
						"id": id,
					})

					return product, findErr
				},
			},

			/* Delete product by id
			http://localhost:8080/graphql?query=mutation+_{delete(id:1){id,name,info,price}}
			*/
			"delete": &graphql.Field{
				Type:        ProductType,
				Description: "Delete product by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					// Lookup existing product
					product, findErr := findProduct(dialect, db, goqu.Ex{
						"id": id,
					})
					if findErr != nil {
						return nil, findErr
					}

					// Remove the existing product
					deleteDialect := goqu.From("products").Where(goqu.Ex{
						"id": id,
					}).Delete()
					deleteQuery, _, toSQLErr := deleteDialect.ToSQL()
					if toSQLErr != nil {
						return nil, toSQLErr
					}

					deleteRes, deleteErr := db.Query(deleteQuery)
					if deleteErr != nil {
						return nil, deleteErr
					}

					defer deleteRes.Close()

					return product, nil
				},
			},
		},
	})
}
