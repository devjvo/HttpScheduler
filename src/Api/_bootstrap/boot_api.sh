#!/bin/sh

# live reload
air -c /app/src/Api/_bootstrap/.air.toml

# start go HTTP server
go run /app/src/Api/main.go
