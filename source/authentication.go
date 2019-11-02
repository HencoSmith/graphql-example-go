package source

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/doug-martin/goqu/v8"

	"github.com/HencoSmith/graphql-example-go/models"
)

// GetUser - Lookup the user based on context returns the user model or alternatively
// an error
func GetUser(currentContext context.Context, dialect goqu.DialectWrapper, db *sql.DB) (models.User, error) {
	// Extract the token from the header
	contextValue := currentContext.Value(models.ContextKey{Key: "header"}).(http.Header)
	authorizationToken := contextValue.Get("Authorization")

	// Check if the token is valid
	// TODO: JWT DB Lookup
	// TODO: JWT validation
	// TODO: JWT user ID extraction
	if authorizationToken != "test" {
		return models.User{}, errors.New("Invalid Token")
	}
	userID := "d56d4bff-4e7e-4cf9-a3d2-38973c9dd57d"

	// Find user
	dialectString := dialect.From("users").Where(goqu.Ex{
		"id":         userID,
		"deleted_at": nil,
	})
	query, _, dialectErr := dialectString.ToSQL()
	if dialectErr != nil {
		return models.User{}, dialectErr
	}

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		return models.User{}, queryErr
	}
	defer rows.Close()

	var usersArr []models.User
	for rows.Next() {
		var row = models.User{}
		scanErr := rows.Scan(
			&row.ID,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
			&row.Email,
		)
		if scanErr != nil {
			return models.User{}, scanErr
		}
		usersArr = append(usersArr, row)
	}
	if errRows := rows.Err(); errRows != nil {
		return models.User{}, errRows
	}

	if len(usersArr) < 1 {
		return models.User{}, errors.New("User Not Found")
	}

	return usersArr[0], nil
}
