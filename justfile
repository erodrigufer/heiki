set shell := ["/bin/sh", "-c"]

DBMATE_GENERAL_OPTIONS := "--migrations-dir ./db/migrations"

default:
  @just --list

# go vet.
[group('go')]
vet:
  cd backend && go vet ./...

# Start the server with air.
[group('go')]
dev: migrate
  cd backend && air

# generate templ files.
[group('go')]
templ:
  templ generate -path ./backend/internal/views

# remove temporary files.
clean:
  cd backend && rm -rf ./tmp

# build for deployment.
[group('deployment')]
build:
  cd backend && env GOOS=freebsd GOARCH=amd64 go build -o ../build/serenitynow ./cmd/serenitynow

# open sqlite cli.
[group('sql')]
sql:
  cd db && sqlite3 database.sqlite3

# run database migrations.
[group('sql')]
migrate: 
    dbmate {{DBMATE_GENERAL_OPTIONS}} up

# rollback the last database migration.
[group('sql')]
rollback: 
    dbmate {{DBMATE_GENERAL_OPTIONS}} rollback

# drop the db managed by dbmate.
[group('sql')]
drop:
    dbmate {{DBMATE_GENERAL_OPTIONS}} drop

# create a new migration file with name `MIGRATION_NAME`.
[group('sql')]
new MIGRATION_NAME:
  dbmate {{DBMATE_GENERAL_OPTIONS}} new {{MIGRATION_NAME}}
