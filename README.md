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

# TODOS
* Test Cases
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