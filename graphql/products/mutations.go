package products

import (
	"math/rand"
	"time"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"
)

func Mutations(products *[]models.Product) *graphql.Object {
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
					rand.Seed(time.Now().UnixNano())
					product := models.Product{
						ID:    int64(rand.Intn(100000)), // generate random ID
						Name:  params.Args["name"].(string),
						Info:  params.Args["info"].(string),
						Price: params.Args["price"].(float64),
					}
					*products = append(*products, product)
					return product, nil
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
					product := models.Product{}
					for i, p := range *products {
						if int64(id) == p.ID {
							if nameOK {
								(*products)[i].Name = name
							}
							if infoOK {
								(*products)[i].Info = info
							}
							if priceOK {
								(*products)[i].Price = price
							}
							product = (*products)[i]
							break
						}
					}
					return product, nil
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
					product := models.Product{}
					for i, p := range *products {
						if int64(id) == p.ID {
							product = (*products)[i]
							// Remove from product list
							*products = append((*products)[:i], (*products)[i+1:]...)
						}
					}

					return product, nil
				},
			},
		},
	})
}
