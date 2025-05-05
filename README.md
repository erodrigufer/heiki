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
# ---- Application ----
export ENVIRONMENT="local"
export PORT=":80"
export SQLITE_PATH="<path to local db for dev>"
export DISABLE_AUTH="<true | false>"
export AUTH_USERNAME
export AUTH_PASSWORD
# ---- dbmate ----
export DATABASE_URL="sqlite:db/database.sqlite3"
```
