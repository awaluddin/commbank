# commbank

## Staircase

### Run

```sh
./staircase -num=5
```
fill number of stair into attribute `-num`


## CRUD

### Database
please provide database `commbank` with single table `employee`
or import from `commbank.sql`

### Configuration
please set config database as your database local
```sh
export DB_USER="root"
export DB_PASS="example"
export DB_NAME="commbank"
export DB_PORT="3307"
export DB_HOST="localhost"
```

### Run

```sh
./crud 
```
crud service will running on http://localhost:8080

### Endpoints

|Method|URL|Description|
|---------|-----------|--------|
|GET | /list | get list eployee data |
|GET | /detail?id= | get detail employee by id |
|POST | /insert | submit new employee|
|PUT | /update | update employee data|
|DELETE | /delete?id= | delte one employee data |
