REPO_ROOT?=$(shell pwd)

# App

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./... -v -race -cover

.PHONY: mock
mock:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run

# Db

.PHONY: demodbs
demodbs:
	docker run -d --name dbui-mysql-demo -p 3316:3306 -e MYSQL_ROOT_PASSWORD=demo genschsa/mysql-employees
	docker run -d --name dbui-postgresql-demo -p 5432:5432 ghusta/postgres-world-db:2.4-alpine

.PHONY: demodbs/destroy
demodbs/destroy:
	docker stop dbui-mysql-demo
	docker stop dbui-postgresql-demo
	docker rm dbui-mysql-demo
	docker rm dbui-postgresql-demo
