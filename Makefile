lint:
	golangci-lint-1.39.0 run -c .golangci.yml ./...

test:
	go test -count=1 -v ./test

bench:
	go test -bench . | tee ./report/report.txt

report:
	./report.py

.PHONY: report