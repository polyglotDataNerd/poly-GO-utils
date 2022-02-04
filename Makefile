# Go Parameters
GOCMD=go
GOTEST=$(GOCMD) test


gotest:
	$(GOTEST) ./test -c -o tests
	./tests
	rm -r ./tests

clean:
	$(GOCMD) clean
