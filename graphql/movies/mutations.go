package movies

import (
	"database/sql"
	"time"

	"github.com/HencoSmith/graphql-example-go/models"
	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
	uuid "github.com/satori/go.uuid"
)

func findMovie(dialect goqu.DialectWrapper, db *sql.DB, expression goqu.Ex) (*models.Movie, error) {
	// Find movie inserted (ID is generated by DB hence lookup is required)
	dialectString := dialect.From("movies").Where(expression)
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

// Mutations - all GraphQL mutations related to movies
func Mutations(dialect goqu.DialectWrapper, db *sql.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type:        MovieType,
				Description: "Create new movie",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"releaseYear": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Insert the new movie
					name, _ := params.Args["name"].(string)
					description, _ := params.Args["description"].(string)
					releaseYear, _ := params.Args["releaseYear"].(int)
					insertDialect := goqu.Insert("movies").Rows(
						goqu.Record{
							"id":           uuid.NewV4(),
							"name":         name,
							"description":  description,
							"release_year": releaseYear,
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

					movie, findErr := findMovie(dialect, db, goqu.Ex{
						// Name has a unique constraint
						"name": params.Args["name"].(string),
					})

					return movie, findErr
				},
			},

			"update": &graphql.Field{
				Type:        MovieType,
				Description: "Update movie by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"releaseYear": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)

					// Update the existing movie
					updateDialect := goqu.Update("movies").Set(
						goqu.Record{
							"name":         params.Args["name"].(string),
							"description":  params.Args["description"].(string),
							"release_year": params.Args["releaseYear"].(int),
						},
					).Where(goqu.Ex{
						"id":         id,
						"deleted_at": nil,
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

					movie, findErr := findMovie(dialect, db, goqu.Ex{
						"id": id,
					})

					return movie, findErr
				},
			},

			"delete": &graphql.Field{
				Type:        MovieType,
				Description: "Delete movie by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)

					// Lookup existing movie
					movie, findErr := findMovie(dialect, db, goqu.Ex{
						"id": id,
					})
					if findErr != nil {
						return nil, findErr
					}

					// Remove the existing movie
					deleteDialect := goqu.Update("movies").Set(
						goqu.Record{
							"deleted_at": time.Now().Format(time.RFC3339),
						},
					).Where(goqu.Ex{
						"id": id,
					})
					deleteQuery, _, toSQLErr := deleteDialect.ToSQL()
					if toSQLErr != nil {
						return nil, toSQLErr
					}

					deleteRes, deleteErr := db.Query(deleteQuery)
					if deleteErr != nil {
						return nil, deleteErr
					}

					defer deleteRes.Close()

					return movie, nil
				},
			},
		},
	})
}
