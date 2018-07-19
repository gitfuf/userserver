# userserver
Test REST API server for user managment. 
Was written in order to practice work with RESP API and SQL databases like Postgresql, MySQL or MSSQL.
In order to run server:
go run ./cmd/server/main.go

Configuration is made through config/config.yaml:

### DB configuration
"db_driver": "postgres" if want use Postgresql
"db_driver": "mssql" if want use MSSQL
"db_driver": "mysql" if want use MySQL
Each SQL database contains pro and test configurations.

Test configuration is used for go test command. Using docker-compose.pg.mysql.mssql.yml allow to run all three databases in order to check userserver work. Steps for running:
1) First need to create userserver image using Dockerfile
docker-compose -f docker-compose.pg.mysql.mssql.yml build 
2) Then start all
docker-compose -f docker-compose.pg.mysql.mssql.yml up
3)To finish and correct stop all use
docker-compose -f docker-compose.pg.mysql.mssql.yml down

### HTTP configuration
"http_port": "8080" declare HTTP port want to use for HTTP requests

### REST API routes:
used Gorrila Mux library (github.com/gorilla/mux)

1)Add new user: "/user"
curl -X POST http://localhost:8080/user -d '{"age":44,"first_name":"Mark","last_name":"Salt","email":"fuf@fu1.com"}'

2)Get user info: "/user/{id:[0-9]+}" 
curl -X GET http://localhost:8080/user/9

3)Update user info: "/user/{id:[0-9]+}" 
curl http://localhost:8080/user/1 -X PUT -d '{"age":24,"first_name":"Maria","last_name":"Solo","email":"ku@fu3.com"}'

4)Delete user: "/user/{id:[0-9]+}" 
curl -X DELETE http://localhost:8080/user/9

5)Also added  "/shutdown" route for remote graceful shutdown. 
curl -X GET http://localhost:8080/shutdown


### Table model
For work is used only one simple table 'users'. Model lookes like: 
`type User struct {`
`	ID        int64  `json:"id"``
`	Age       int    `json:"age"``
`	FirstName string `json:"first_name"``
`	LastName  string `json:"last_name"``
`	Email     string `json:"email"``
`}`