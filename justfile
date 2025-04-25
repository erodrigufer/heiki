set shell := ["/bin/zsh", "-c"]

default:
  @just --list

# go vet.
[group('go')]
vet:
  cd backend && go vet ./...

# Start the server with air.
[group('go')]
dev:
  cd backend && air

# generate templ files.
[group('go')]
templ:
  templ generate -path ./backend/internal/views

# remove temporary files.
clean:
  cd backend && rm -rf ./tmp
