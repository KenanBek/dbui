REPO_ROOT?=$(shell pwd)

# App

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./... -race -cover -tags=component

.PHONY: mock
mock:
	go generate ./...

# Db

.PHONY: demodbs
demodbs:
	docker run -d --name mysql-employees -p 3316:3306 -e MYSQL_ROOT_PASSWORD=demo -v $(REPO_ROOT)/demodata/mysqlemployees:/var/lib/mysql genschsa/mysql-employees
	docker run -d --name postgresql-world-db -p 5432:5432 ghusta/postgres-world-db:2.5

.PHONY: demodbs/destroy
demodbs/destroy:
	docker stop mysql-employees
	docker stop postgresql-world-db
	docker rm mysql-employees
	docker rm postgresql-world-db
