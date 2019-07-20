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

# Database Setup
```
docker pull postgres
docker run --name postgres-container -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=test_db -d postgres
```

# TODOS
* Subscriptions
* DB (containerized maybe?)
* Code Coverage
* Documentation Generation
* Performance Testing
* API Analytics
* Packaging
* Workers
* Graceful Shutdown
* heroku app hosting
* JWT