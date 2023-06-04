.PHONY: db-up
db-up:
	docker compose up -d db-admin

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