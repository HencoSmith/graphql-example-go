package users

import (
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v8"
	"github.com/graphql-go/graphql"

	source "github.com/HencoSmith/graphql-example-go/source"
)

// Queries - all GraphQL queries related to movies
func Queries(dialect goqu.DialectWrapper, db *sql.DB) graphql.Fields {
	return graphql.Fields{
		"getToken": &graphql.Field{
			Type:        graphql.String,
			Description: "Return a JWT for the specified user",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email, _ := p.Args["email"].(string)
				password, _ := p.Args["password"].(string)

				// Lookup encrypted password from DB
				user, err := source.GetUser(dialect, db, "", email)
				if err != nil {
					return nil, err
				}

				valid, err := source.ValidHash(password, user.EncryptedPassword)
				if err != nil {
					return nil, err
				}

				if !valid {
					return nil, errors.New("User Not Found")
				}

				// TODO: generate new JWT here
				return "test", nil
			},
		},
	}
}
