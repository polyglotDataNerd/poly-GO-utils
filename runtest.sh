#!/usr/bin/env bash
go test ./test -c -o tests
./tests
rm -r ./tests
