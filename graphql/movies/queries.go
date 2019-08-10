package movies

import (
	"database/sql"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
)

// Queries - all GraphQL queries related to movies
func Queries(dialect goqu.DialectWrapper, db *sql.DB) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"movie": &graphql.Field{
					Type:        MovieType,
					Description: "Get movie by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)
						if ok {
							// Find movie
							dialectString := dialect.From("movies").Where(goqu.Ex{
								"id":         id,
								"deleted_at": nil,
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

							var moviesArr []models.Movie
							for rows.Next() {
								var movieRow = models.Movie{}
								scanErr := rows.Scan(
									&movieRow.ID,
									&movieRow.CreatedAt,
									&movieRow.UpdatedAt,
									&movieRow.DeletedAt,
									&movieRow.Name,
									&movieRow.ReleaseYear,
									&movieRow.Description,
									&movieRow.Rating,
									&movieRow.ReviewCount,
								)
								if scanErr != nil {
									return nil, scanErr
								}
								moviesArr = append(moviesArr, movieRow)
							}
							if errRows := rows.Err(); errRows != nil {
								return nil, errRows
							}

							if len(moviesArr) < 1 {
								return nil, nil
							}

							return &moviesArr[0], nil
						}
						return nil, nil
					},
				},

				"list": &graphql.Field{
					Type:        graphql.NewList(MovieType),
					Description: "Get movie list",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						dialectString := dialect.From("movies").Where(goqu.Ex{
							"deleted_at": nil,
						}).Order(goqu.C("id").Asc())
						query, _, dialectErr := dialectString.ToSQL()
						if dialectErr != nil {
							return nil, dialectErr
						}

						rows, queryErr := db.Query(query)
						if queryErr != nil {
							return nil, queryErr
						}
						defer rows.Close()

						var moviesArr []models.Movie
						for rows.Next() {
							var movieRow = models.Movie{}
							scanErr := rows.Scan(
								&movieRow.ID,
								&movieRow.CreatedAt,
								&movieRow.UpdatedAt,
								&movieRow.DeletedAt,
								&movieRow.Name,
								&movieRow.ReleaseYear,
								&movieRow.Description,
								&movieRow.Rating,
								&movieRow.ReviewCount,
							)
							if scanErr != nil {
								return nil, scanErr
							}
							moviesArr = append(moviesArr, movieRow)
						}
						if errRows := rows.Err(); errRows != nil {
							return nil, errRows
						}

						return &moviesArr, nil
					},
				},
			},
		},
	)
}
