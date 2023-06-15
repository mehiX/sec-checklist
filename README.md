# Security checks REST API

The security checks needed for auditing are initially provided in an Excel spreadsheet. The aim of the project is to load this data into a relational database and provide a REST API to query these controls. 

## Run with docker compose

Needs: 
- `docker compose`

Create the environment file (replace the values with the correct ones for your environment):

```shell
cat > .env <<EOF
EXCEL_FILEPATH=/data/ISMS1042_VIT_v0.04.xlsx
SHEET_NAME='ISMS1042 6.2 with all labels'
EOF
```

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

./dist/secctrls api load --from <excel_file_path> --fromSheet <sheet_name>
```

Start the webserver:

```shell
# on Windows use ./dist/secctrls.exe
./dist/secctrls api serve --http 127.0.0.1:8080 --db
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
curl -s http://localhost:8080/docs/controls/filter | jq '.'
```
```json
{
  "Request body example": {
    "only_handle_centrally": true,
    "handled_centrally_by": "BSO",
    "exclude_for_external_supplier": true,
    "software_development_relevant": true,
    "cloud_only": true,
    "physical_security_only": true,
    "personal_security_only": true
  },
  "Response example": [
    {
      "Type": "",
      "ID": "",
      "Name": "",
      "Description": "",
      "C": "",
      "I": "",
      "A": "",
      "T": "",
      "PD": "",
      "NSI": "",
      "SESE": "",
      "OTCL": "",
      "CSRDirection": "",
      "SPSA": "",
      "SPSAUnique": "",
      "GDPR": false,
      "GDPRUnique": false,
      "ExternalSupplier": false,
      "AssetType": "",
      "OperationalCapability": "",
      "PartOfGISR": false,
      "LastUpdated": "",
      "OldID": "",
      "OnlyHandledCentrally": false,
      "HandledCentrallyBy": "",
      "ExcludeForExternalSupplier": false,
      "SoftwareDevelopmentRelevant": false,
      "CloudOnly": false,
      "PhysicalSecurityOnly": false,
      "PersonalSecurityOnly": false
    }
  ]
}
```

## TODO

- create a readonly Mysql user and use that one for normal operations. A user with create rights is only needed to initialize the database
