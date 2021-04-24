lint:
	golangci-lint-1.39.0 run -c .golangci.yml ./...

test:
	go test -count=1 -v ./test

bench:
	go test -bench . > ./report/report.txt

.PHONY: test