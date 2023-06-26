#!/bin/bash

go build -o ./bin/bookings cmd/web/*.go

./bin/bookings -dbname=booking -dbuser=booker -dbpass=qwerty1! -production=false -cache=false