package source

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/template"

	configStruct "github.com/HencoSmith/graphql-example-go/config/struct"

	"github.com/doug-martin/goqu/v8"
)

func getConnectionString(config configStruct.Configuration) (string, error) {
	// Lookup details defined in environment variables and overwrite config values
	password := config.Database.Password
	envPassword := os.Getenv("POSTGRES_PASSWORD")
	if len(envPassword) != 0 {
		password = envPassword
	}
	user := config.Database.User
	envUser := os.Getenv("POSTGRES_USER")
	if len(envUser) != 0 {
		user = envUser
	}
	DBName := config.Database.Name
	envDBName := os.Getenv("POSTGRES_DB")
	if len(envDBName) != 0 {
		DBName = envDBName
	}

	input := configStruct.DatabaseConfiguration{
		Name:     DBName,
		User:     user,
		Password: password,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		SSL:      config.Database.SSL,
	}

	connectionString := `postgresql://{{.User}}:{{.Password}}@{{.Host}}:{{.Port}}/{{.Name}}?sslmode={{.SSL}}`
	templateString := template.Must(template.New("connectionString").Parse(connectionString))
	var stringParsed bytes.Buffer
	if errExecute := templateString.Execute(&stringParsed, input); errExecute != nil {
		return "", errExecute
	}

	return stringParsed.String(), nil
}

// ConnectToDB attempts to connect to the database and returns a pointer to the database
// along with an error if applicable
func ConnectToDB(config configStruct.Configuration) (*sql.DB, error) {
	connStr, err := getConnectionString(config)
	if err != nil {
		return nil, err
	}

	fmt.Println("connecting to: " + connStr)
	return sql.Open("postgres", connStr)
}

func seedDB(db *sql.DB, tableName string, records []interface{}) error {
	insertDialect := goqu.Insert(tableName).Rows(records...)
	insertQuery, _, SQLErr := insertDialect.ToSQL()
	if SQLErr != nil {
		return SQLErr
	}

	insertRes, insertErr := db.Query(insertQuery)
	// Ignore constraint errors, which means that the data has already been inserted
	if insertErr != nil {
		if !strings.Contains(insertErr.Error(), `unique constraint "`+tableName+`_pkey"`) {
			return insertErr
		}
	} else {
		defer insertRes.Close()
	}
	return nil
}

// InitTables - Create and load DB tables with data
// returns an error or nil if no error ocurred
func InitTables(dialect goqu.DialectWrapper, db *sql.DB) error {
	fmt.Println("initializing DB...")
	// Create the table(s) if it does not already exist
	createTable, createErr := db.Query(`
	CREATE TABLE IF NOT EXISTS public.movies
	(
		id uuid NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT now(),
		updated_at timestamp with time zone NOT NULL DEFAULT now(),
		deleted_at timestamp with time zone,
		users_id uuid NOT NULL,
		name character varying(128) NOT NULL,
		release_year integer NOT NULL,
		description text,
		rating numeric NOT NULL DEFAULT '0.0',
		review_count bigint NOT NULL DEFAULT 0,
		PRIMARY KEY (id)
	)
	WITH (
		OIDS = FALSE
	);

	ALTER TABLE public.movies
		OWNER to "user";

	CREATE TABLE IF NOT EXISTS public.movies_reviews
	(
		id uuid NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT now(),
		updated_at timestamp with time zone NOT NULL DEFAULT now(),
		deleted_at timestamp with time zone,
		movies_id uuid NOT NULL,
		users_id uuid NOT NULL,
		rating numeric NOT NULL,
		PRIMARY KEY (id)
	)
	WITH (
		OIDS = FALSE
	);

	ALTER TABLE public.movies_reviews
		OWNER to "user";

	CREATE TABLE IF NOT EXISTS public.users
	(
		id uuid NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT now(),
		updated_at timestamp with time zone NOT NULL DEFAULT now(),
		deleted_at timestamp with time zone,
		email character varying(64) NOT NULL,
		encrypted_password character varying(512) NOT NULL,
		PRIMARY KEY (id)
	)
	WITH (
		OIDS = FALSE
	)
	TABLESPACE pg_default;
	
	ALTER TABLE public.users
		OWNER to "user";

	DROP INDEX IF EXISTS movies_id_idx;

	CREATE INDEX movies_id_idx
		ON public.movies USING btree
		(id ASC NULLS LAST)
		TABLESPACE pg_default;

	DROP INDEX IF EXISTS movies_deleted_at_idx;

	CREATE INDEX movies_deleted_at_idx
		ON public.movies USING btree
		(deleted_at ASC NULLS FIRST)
		TABLESPACE pg_default;

	ALTER TABLE public.movies
		DROP CONSTRAINT IF EXISTS movies_users_id_fkey;
	
	ALTER TABLE public.movies
		ADD CONSTRAINT movies_users_id_fkey FOREIGN KEY (users_id)
		REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION;

	DROP INDEX IF EXISTS fki_movies_users_id_fkey;
	
	CREATE INDEX fki_movies_users_id_fkey
		ON public.movies(users_id);

	DROP INDEX IF EXISTS movies_reviews_id_idx;

	CREATE INDEX movies_reviews_id_idx
		ON public.movies_reviews USING btree
		(id ASC NULLS LAST)
		TABLESPACE pg_default;

	DROP INDEX IF EXISTS movies_reviews_deleted_at_idx;

	CREATE INDEX movies_reviews_deleted_at_idx
		ON public.movies_reviews USING btree
		(deleted_at ASC NULLS FIRST)
		TABLESPACE pg_default;

	ALTER TABLE public.movies_reviews
		DROP CONSTRAINT IF EXISTS movies_reviews_movies_id_fkey;

	ALTER TABLE public.movies_reviews
		ADD CONSTRAINT movies_reviews_movies_id_fkey FOREIGN KEY (movies_id)
		REFERENCES public.movies (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION;

	DROP INDEX IF EXISTS fki_movies_reviews_movies_id_fkey;
	
	CREATE INDEX fki_movies_reviews_movies_id_fkey
		ON public.movies_reviews(movies_id);

	ALTER TABLE public.movies_reviews
		DROP CONSTRAINT IF EXISTS movies_reviews_users_id_fkey;

	ALTER TABLE public.movies_reviews
		ADD CONSTRAINT movies_reviews_users_id_fkey FOREIGN KEY (users_id)
		REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE NO ACTION
		ON DELETE NO ACTION;
	
	DROP INDEX IF EXISTS fki_movies_reviews_users_id_fkey;

	CREATE INDEX fki_movies_reviews_users_id_fkey
		ON public.movies_reviews(users_id);

	DROP INDEX IF EXISTS users_id_idx;

	CREATE INDEX users_id_idx
		ON public.users USING btree
		(id ASC NULLS LAST)
		TABLESPACE pg_default;

	DROP INDEX IF EXISTS users_deleted_at_idx;

	CREATE INDEX users_deleted_at_idx
		ON public.users USING btree
		(deleted_at ASC NULLS FIRST)
		TABLESPACE pg_default;
	`)
	if createErr != nil {
		return createErr
	}
	defer createTable.Close()

	// Hash a default password
	encryptedPassword, hashErr := Hash("test")
	if hashErr != nil {
		return hashErr
	}

	// Build Users Table Seed
	usersSeedErr := seedDB(db, "users", []interface{}{
		goqu.Record{
			"id":                 "d56d4bff-4e7e-4cf9-a3d2-38973c9dd57d",
			"email":              "test@mail.com",
			"encrypted_password": encryptedPassword,
		},
	})
	if usersSeedErr != nil {
		return usersSeedErr
	}

	// Build Movies Table Seed
	movieSeedErr := seedDB(db, "movies", []interface{}{
		goqu.Record{
			"id":           "13cbd25a-4a9d-4e71-9c39-4fc515083c95",
			"name":         "Scary Stories to Tell in the Dark",
			"release_year": 2019,
			"description":  "A group of teens face their fears in order to save their lives.",
			"users_id":     "d56d4bff-4e7e-4cf9-a3d2-38973c9dd57d",
		},
		goqu.Record{
			"id":           "77034dd5-d3e4-4a44-a7fa-c2730dfe5370",
			"name":         "Dora and the Lost City of Gold",
			"release_year": 2019,
			"description":  "Dora, a teenage explorer, leads her friends on an adventure to save her parents and solve the mystery behind a lost city of gold.",
			"users_id":     "d56d4bff-4e7e-4cf9-a3d2-38973c9dd57d",
		},
		goqu.Record{
			"id":           "a774e5ff-a5f9-4643-832d-27d131344fe3",
			"name":         "The Art of Racing in the Rain",
			"release_year": 2019,
			"description":  "Through his bond with his owner, aspiring Formula One race car driver Denny, golden retriever Enzo learns that the techniques needed on the racetrack can also be used to successfully navigate the journey of life.",
			"users_id":     "d56d4bff-4e7e-4cf9-a3d2-38973c9dd57d",
		},
	})
	if movieSeedErr != nil {
		return movieSeedErr
	}

	fmt.Println("OK")
	return nil
}
