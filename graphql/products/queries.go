package products

// TODO: rely on DB / File to hook up products

// import (
// 	"github.com/graphql-go/graphql"
// )

// var queryType = graphql.NewObject(
// 	graphql.ObjectConfig{
// 		Name: "Query",
// 		Fields: graphql.Fields{
// 			/* Get (read single product by id)
// 			http://localhost:8080/product?query={product(id:1){name,info,price}}
// 			*/
// 			"product": &graphql.Field{
// 				Type:        ProductType,
// 				Description: "Get product by id",
// 				Args: graphql.FieldConfigArgument{
// 					"id": &graphql.ArgumentConfig{
// 						Type: graphql.Int,
// 					},
// 				},
// 				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 					id, ok := p.Args["id"].(int)
// 					if ok {
// 						// Find product
// 						for _, product := range products {
// 							if int(product.ID) == id {
// 								return product, nil
// 							}
// 						}
// 					}
// 					return nil, nil
// 				},
// 			},
// 			/* Get (read) product list
// 			http://localhost:8080/product?query={list{id,name,info,price}}
// 			*/
// 			"list": &graphql.Field{
// 				Type:        graphql.NewList(ProductType),
// 				Description: "Get product list",
// 				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
// 					return products, nil
// 				},
// 			},
// 		},
// 	},
// )
