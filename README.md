### firt, set your variale
export POSTGRESQL_URL='postgres://user:password@localhost:5432/tan-test?sslmode=disable'
### if you want to up
migrate -database ${POSTGRESQL_URL} -path ./migration up

### if you want to down
migrate -database ${POSTGRESQL_URL} -path ./migration down