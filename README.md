# serenity now

平気

## Requirements

- just
- air
- templ
- go
- direnv

## Commands

- Run `just` to see a list of all `just` targets.

## Environment variables

```bash
# ---- Application ---
export ENVIRONMENT="local"
export PORT=":80"
export DATABASE_URL="sqlite:db/database.sqlite3" # required for dbmate
export SQLITE_PATH="<path to local db for dev>"
```
