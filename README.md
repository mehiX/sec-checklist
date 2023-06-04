# Security checks REST API

The security checks needed for auditing are initially provided in an Excel spreadsheet. The aim of the project is to load this data into a relational database and provide a REST API to query these controls. 

## Run with docker compose

Needs: 
- `docker compose`

```shell
make up
```

This starts:

- the database (mariadb)
- a database client (adminer)
- the application 

## Build and run the binary

Needs:

- go >= 1.18
- `docker compose` to start the database or access to a MariaDB/Mysql instance

Create `.env` by providing your database connection credentials:

```shell
cat > .env <<EOF
CHECKLISTS_DSN=test:test@tcp(127.0.0.1)/test
EOF
```

Build the binary:

```shell
make binary
```

Load initial data for the Excel spreadsheet:

```shell
# on Windows use ./dist/secctrls.exe

./dist/secctrls -from <excel_file_path> -fromSheet <sheet_name> -db -init
```

Start the webserver:

```shell
# on Windows use ./dist/secctrls.exe
./dist/secctrls -db -http 127.0.0.1:8080
```

## Endpoints

List all loaded controls:

```shell
curl http://localhost:8080/controls/
```

Get details about a specific control:

```shell
curl http://localhost:8080/controls/5.1.2.2
```

Filter controls based on the application profile:

```shell
curl \
    -d '{"only_handle_centrally": true}' \
    -H "Content-type: application/json" \
    http://127.0.0.1:8080/controls | jq '.'
```

Get help for the available filters:

```shell
curl http://localhost:8080/docs/controls/filter
```

## TODO

- create a readonly Mysql user and use that one for normal operations. A user with create rights is only needed to initialize the database
