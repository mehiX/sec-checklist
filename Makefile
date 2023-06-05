EXCEL_FILE=./data/ISMS1042_VIT_v0.04.xlsx
EXCEL_SHEET='ISMS1042 6.2 with all labels'

.PHONY: db-up
db-up:
	docker compose up -d db-admin

.PHONY: db-init
db-init: down db-up binary
	./dist/secctrls -from ${EXCEL_FILE} -fromSheet ${EXCEL_SHEET} -db -init

.PHONY: down
down:
	docker compose down --volumes --remove-orphans

.PHONY: up
up: db-up
	docker compose up -d --build secchecklist

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

.PHONY: start
start: binary
	./dist/secctrls -db -http 127.0.0.1:8080