# graphql-example-go
Example of GraphQL implementation in Golang

# Playground
* run the server
* navigate to:
  * Prisma (Preferred) [http://localhost:8080/playground]
  * GraphiQL [http://localhost:8080/graphql]

# Documentation
Refer to playground generated docs for API documentation.
For Golang related documentation:
```
godoc -http=:6060
```

# Configuration
Found in ./config/config.yml
* server - Server related configuration
  * port - HTTP port to host the server on

# Testing
Test cases found in ./test
Startup the server then run:
```
cd test
go test
go test -run TestMyFunc
```

# Environment Variables
These are optional, see configuration file for defaults.
The passwords should set in the environment variables however for this example 
defaults are provided in the configuration file to make set up easier.
```
POSTGRES_PASSWORD - Database.Password
POSTGRES_USER - Database.User
POSTGRES_DB - Database.Name
JWT_KEY - JWT.Key
```

# Database Setup
```
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
* API keys - JWT
* Update test cases once JWT is done
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