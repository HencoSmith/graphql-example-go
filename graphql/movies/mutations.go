package movies

import (
	"database/sql"
	"math"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/doug-martin/goqu/v8"
	uuid "github.com/satori/go.uuid"

	"github.com/HencoSmith/graphql-example-go/models"
	source "github.com/HencoSmith/graphql-example-go/source"
)

// findMovie - lookup a movie matching the specified expression
// dialect - Query builder dialect object used
// db - SQL DB connection to use
// expression - Expression movie looking up should adhere to
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
			&movieRow.UsersID,
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

// calculateMovieRating - Determine the overall sum of ratings for a movie
// dialect - Query builder dialect object used
// db - SQL DB connection to use
// id - UUID of the movie to calculate the rating sum for
func calculateMovieRating(dialect goqu.DialectWrapper, db *sql.DB, id string) (float64, error) {
	// Find movie inserted (ID is generated by DB hence lookup is required)
	dialectString := dialect.From("movies_reviews").Select(goqu.SUM("rating").As("totalRating")).Where(goqu.Ex{
		"movies_id": id,
	})
	query, _, dialectErr := dialectString.ToSQL()
	if dialectErr != nil {
		return -1.0, dialectErr
	}

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		return -1.0, queryErr
	}
	defer rows.Close()

	var ratingsArr []models.RatingSum
	for rows.Next() {
		var ratingRow = models.RatingSum{}
		scanErr := rows.Scan(
			&ratingRow.Total,
		)
		if scanErr != nil {
			return -1.0, scanErr
		}
		ratingsArr = append(ratingsArr, ratingRow)
	}
	if errRows := rows.Err(); errRows != nil {
		return -1.0, errRows
	}

	if len(ratingsArr) < 1 {
		return -1.0, nil
	}

	return ratingsArr[0].Total, nil
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
					user, customError := source.GetUser(params.Context, dialect, db)
					if customError != nil {
						return nil, customError
					}

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
							"users_id":     user.ID,
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
						Type: graphql.String,
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"releaseYear": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user, customError := source.GetUser(params.Context, dialect, db)
					if customError != nil {
						return nil, customError
					}

					id, _ := params.Args["id"].(string)
					name, _ := params.Args["name"].(string)
					description, _ := params.Args["description"].(string)
					releaseYear, _ := params.Args["releaseYear"].(int)
					updatedAt := time.Now().Format(time.RFC3339)

					updateFields := goqu.Record{}
					if len(name) > 0 {
						updateFields["name"] = name
					}
					if len(description) > 0 {
						updateFields["description"] = description
					}
					if releaseYear > 1900 {
						updateFields["release_year"] = releaseYear
					}
					updateFields["updated_at"] = updatedAt

					// Update the existing movie
					updateDialect := goqu.Update("movies").Set(
						updateFields,
					).Where(goqu.Ex{
						"id":         id,
						"deleted_at": nil,
						"users_id":   user.ID,
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
					user, customError := source.GetUser(params.Context, dialect, db)
					if customError != nil {
						return nil, customError
					}

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
						"id":       id,
						"users_id": user.ID,
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

			"rate": &graphql.Field{
				Type:        graphql.String,
				Description: "Rate a movie by ID. Returns 'success' / 'failure'",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"rating": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "0 - 10 (0 - worst; 10 - best)",
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user, customError := source.GetUser(params.Context, dialect, db)
					if customError != nil {
						return nil, customError
					}

					id, _ := params.Args["id"].(string)
					rating, _ := params.Args["rating"].(int)

					// Limit rating 0 - 10
					formattedRating := math.Max(float64(rating), float64(0))
					formattedRating = math.Min(float64(formattedRating), float64(10))

					// Find the current rating
					movie, findErr := findMovie(dialect, db, goqu.Ex{
						"id": id,
					})

					if findErr != nil {
						return "failure", findErr
					}

					// Add the rating as a running total
					insertDialect := goqu.Insert("movies_reviews").Rows(
						goqu.Record{
							"id":        uuid.NewV4(),
							"rating":    formattedRating,
							"movies_id": movie.ID,
							"users_id":  user.ID,
						},
					)
					insertQuery, _, insertToSQLErr := insertDialect.ToSQL()
					if insertToSQLErr != nil {
						return "failure", insertToSQLErr
					}

					insertRes, insertErr := db.Query(insertQuery)
					if insertErr != nil {
						return "failure", insertErr
					}

					defer insertRes.Close()

					// Calculate review total (aggregate functions not allowed on update)
					totalRating, ratingErr := calculateMovieRating(dialect, db, id)
					if ratingErr != nil {
						return "failure", ratingErr
					}

					// determine total amount of reviewers
					totalReviewers := float64(int32(movie.ReviewCount) + int32(1))

					// Update the existing movie
					updateDialect := goqu.Update("movies").Set(
						goqu.Record{
							"rating":       totalRating / totalReviewers,
							"review_count": totalReviewers,
						},
					).Where(goqu.Ex{
						"movies.id":         id,
						"movies.deleted_at": nil,
					})
					updateQuery, _, toSQLErr := updateDialect.ToSQL()
					if toSQLErr != nil {
						return "failure", toSQLErr
					}

					updateRes, updateErr := db.Query(updateQuery)
					if updateErr != nil {
						return "failure", updateErr
					}

					defer updateRes.Close()

					return "success", nil
				},
			},
		},
	})
}
