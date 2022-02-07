# Go Parameters
GOCMD=go
GOTEST=$(GOCMD) test

clean:
	$(GOCMD) clean

gocompose:
	docker-compose up --build app

gotestdocker:
	docker-compose up --build test
	docker-compose down localstack
	docker-compose down test

