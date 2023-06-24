# Booking App
Go application for booking/reservation system with Monolithic architecture.

- Built in Go version 1.20
- Uses the [chi router](https://github.com/go-chi/chi)
- Uses [Alex Edward's SCS](https://github.com/alexedwards/scs) session management
- Uses the [nosurf](https://www.github.com/justinas/nosurf)

---
## Commands

- Run project (Linux):
````
./run.sh
````

- Run go test:
````
go test -v
````

- Run all go tests in subdirectories:
````
go test -v ./...
````

- Run coverage test command: 
````
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
````

- Run migration tool (parameters: up/down/reset):
````
~/go/bin/soda migrate
````