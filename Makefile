
REPO_ROOT?=$(shell pwd)

# App

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./... -race -tags=component

# Db

.PHONY: demodbs
demodbs:
	docker run -d --name mysql-employees -p 3316:3306 -e MYSQL_ROOT_PASSWORD=demo -v $(REPO_ROOT)/demodata/mysqlemployees:/var/lib/mysql genschsa/mysql-employees

.PHONY: demodbs/destroy
demodbs/destroy:
	docker stop mysql-employees
	docker rm mysql-employees
