### Working with database [Postgres]
#### Backend #2 (Install & use Docker + Postgres + TablePlus to create DB schema)
brew install docker colima
colima start
docker pull postgres:12-alpine
docker run --name postgres12 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pass -d postgres:12-alpine

#### Backend #3 (How to write & run database migration in Golang)
brew install golang-migrate
migrate create -ext sql -dir db/migration -seq init_schema
docker exec -it postgres12 createdb -U root simple_bank
migrate -path db/migration -database "postgresql://root:pass@localhost:5433/simple_bank?sslmode=disable" -verbose up

#### Backend #4 (Generate CRUD Golang code from SQL | Compare db/sql, gorm, sqlx & sqlc)
brew install sqlc
sqlc init

#### Backend #5 (Write Golang unit tests for database CRUD with random data)
#### Backend #6 (A clean way to implement database transaction in Golang)

#### Backend #7 (DB transaction lock & How to handle deadlock in Golang)
https://wiki.postgresql.org/wiki/Lock_Monitoring
https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html
https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-ROWS

#### Backend #8 (How to avoid deadlock in DB transaction? Queries order matters!)