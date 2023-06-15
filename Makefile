EXCEL_FILE=./data/ISMS1042_VIT_v0.04.xlsx
EXCEL_SHEET='ISMS1042 6.2 with all labels'

.PHONY: db-up
db-up:
	docker compose up -d db-admin

.PHONY: db-init
db-init: down db-up binary
	./dist/secctrls api load --from ${EXCEL_FILE} --fromSheet ${EXCEL_SHEET}

.PHONY: down
down:
	docker compose down --volumes --remove-orphans

.PHONY: up
up: db-up
	docker compose up -d --build secchecklist client

dist:
	mkdir -p dist

.PHONY: clean
clean:
	rm -rf dist
	rm -f cover.out

.PHONY: binary
binary:
	./script/make.sh binary

.PHONY: crossbinary-default
crossbinary-default:
	./script/make.sh crossbinary-default

.PHONY: test-unit
test-unit:
	./script/make.sh test-unit

.PHONY: start-api
start-api: binary
	./dist/secctrls api serve --db --http 127.0.0.1:8080

.PHONY: start-client
start-client: binary
	./dist/secctrls client --api http://localhost:8080 --http 127.0.0.1:8081