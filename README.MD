# Go BoilerPlate

## Go Commands

```go
// Initialize a new Go module
go mod init github.com/anuragShingare30/go-boilerplate

// Tidy up the module, adding missing and removing unused modules
go mod tidy
```

## File and Folder Structure

```
backend
    ├── cmd
    │   └── go-boilerplate
    │       └── main.go          # Application entry internal
    │
    ├── internal
    │   ├── config                # config management and loading env variables when server starts
    │   ├── database              # storing DB migrations and connection logic
    │   ├── errs                  # Storing custom errors and error handling logic
    │   ├── handler               # Storing Handler logic
    │   ├── lib                   # Store any other modules or logic
    │   ├── logger                # Storing logger logic and code
    │   ├── middleware            # store global and error handling and request handling middlewares
    │   ├── model                 # storing Go structs which have one to one mapping with DB tables
    │   ├── repository            # Storing DB queries and modules
    │   ├── router                # router related code
    │   ├── server                # spins up the new HTTP server, used in main.go
    │   ├── service               # business logic code
    │   ├── sqlerr                # Code for handling DB errors
    │   ├── testing               # testing related code
    │   ├── validation            # validation related code
```

## Resources

**Project Layout** - https://github.com/golang-standards/project-layout

**echo framework** - "https://echo.labstack.com/docs/quick-start"

**pgx - SQL Driver** - https://github.com/jackc/pgx

**tern - SQL Migrator** - https://github.com/jackc/tern

**zerolog - JSON Logger** - https://github.com/rs/zerolog

**newrelic -Monitoring and Observability** - "https://pkg.go.dev/github.com/newrelic/go-agent/v3@v3.40.1/newrelic"

**validator** - https://github.com/go-playground/validator

**koanf - Configuration Management** - https://github.com/knadh/koanf

**testify - for testing** - https://github.com/stretchr/testify

**taskfile** - https://taskfile.dev/
