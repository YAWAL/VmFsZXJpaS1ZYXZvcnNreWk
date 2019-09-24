# GWP 

## Prerequisites:
1. [Golang](https://golang.org/dl/) (1.12)
2. [Postgres](https://www.postgresql.org/download/) 
3. [Dependency management for Go](https://golang.github.io/dep/docs/installation.html) Dep

## How to start a server:
1. Clone or download repository to your local machine
2. Open ```src``` folder and run command ```dep ensure```
3. Set up config file - provide port and Postgres credential, f.e.:
```
{
    "host": ":8080",
    "database": {
        "dialect": "postgres",
        "user": "yourPostgresUser",
        "db_name": "yourDatabaseName",
        "ssl_mode": "disable",
        "password": "yourPostgresPassword"
    }
}
```
4. Create database in Postgres with the same name provided in config ``` "db_name": "yourDatabaseName"```
5. Navigate to ```src``` directory
6. Run ```go run main.go```