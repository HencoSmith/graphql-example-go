# graphql-example-go
Example of GraphQL implementation in Golang

# Playground
* run the server
* navigate to:
  * Prisma (Preferred) [http://localhost:8080/playground]
  * GraphiQL [http://localhost:8080/graphql]

# Executing API calls
All API calls except 'getToken' requires and Authorization header.
For testing purposes run the following query to obtain a test token:
```
query {
  getToken(email: "test@mail.com", password: "test")
}
```
Then provide this HTTP header in all other API calls e.g.
```
{
  "authorization":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJkNTZkNGJmZi00ZTdlLTRjZjktYTNkMi0zODk3M2M5ZGQ1N2QiLCJleHAiOjE1NzY2NjQ4MjR9.ojFLoZsw1vuXTHValALjwXdiMDOjMZd08Qs6hpEXcFQ"
}
```

# Documentation
Refer to playground generated docs for API documentation.
For Golang related documentation:
```bash
godoc -http=:6060
```

# Configuration
Found in ./config/config.yml
* server - Server related configuration
  * port - HTTP port to host the server on
  * timeout - Cool down before exiting the server, after receiving termination command, in seconds
* database - PostgreSQL DB details
  * user - username
  * host - host address e.g. 'localhost'
  * port - host port e.g. '5432'
  * name - name of the database to connect to
  * password - password associated with user name
  * ssl - SSL mode used during the database connection
* jwt -
  * key - Key used to encrypt JWT tokens with
  * expiration - After how many hours the token should expire

# Testing
Test cases found in ./test
Startup the server then run:
```bash
cd test
go test
```
or alternatively for a specific test case
```bash
go test -run TestGetToken
```

# Environment Variables
These are optional, see configuration file for defaults.
Passwords should always be set in the environment variables however for this example 
defaults are provided in the configuration file to make set up easier.
```
POSTGRES_PASSWORD - Database.Password
POSTGRES_USER - Database.User
POSTGRES_DB - Database.Name
JWT_KEY - JWT.Key
```

# Database Setup
```bash
docker pull postgres
docker run -p 5432:5432 --name postgres-container -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=test_db -d postgres
```

# Authorization
Make the following GraphQL query
```javascript
query {
  getToken(email: "test@mail.com", password: "test")
}
```
Then insert the resulting token in the HTTP Headers of each API call e.g.
```javascript
{
  "authorization": "paste resulting token here..."
}
```

# TODOS
* Display user email on movie queries instead of ID
* Move DB creation string to file
* Subscriptions
* Code Coverage
* Performance Testing
* API Analytics
* Packaging
* Workers
* heroku app hosting
* Update readme with godoc reference, build status, coverage etc. tags