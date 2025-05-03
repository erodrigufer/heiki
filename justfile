set shell := ["/bin/zsh", "-c"]

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

# Run database migrations.
[group('dbmate')]
migrate: 
    dbmate {{DBMATE_GENERAL_OPTIONS}} up

# Rollback the last database migration.
[group('dbmate')]
rollback: 
    dbmate {{DBMATE_GENERAL_OPTIONS}} rollback

# Drop the db managed by dbmate.
[group('dbmate')]
drop:
    dbmate {{DBMATE_GENERAL_OPTIONS}} drop

# Create a new migration file with name `MIGRATION_NAME`.
[group('dbmate')]
new MIGRATION_NAME:
  dbmate {{DBMATE_GENERAL_OPTIONS}} new {{MIGRATION_NAME}}
