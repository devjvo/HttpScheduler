# HttpScheduler

## Requirements
Mage: https://magefile.org/
Go Migrate: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

## Set up
docker build -t custom-go:latest .
go mod tidy

## Containers
### Build and Up
mage init

### Stop and remove containers
mage down

## Database
### Connect
mage DbConnect

### Migrations

mage MigrateCreate table_name

migrate -database postgres://dev:dev@127.0.0.1:2432/http_scheduler?sslmode=disable -path db/migrations up
migrate -database postgres://dev:dev@127.0.0.1:2432/http_scheduler?sslmode=disable -path db/migrations down

git push personal:devjvo/HttpScheduler.git