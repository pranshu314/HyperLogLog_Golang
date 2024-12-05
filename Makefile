.PHONY:
run:
	go run hyperloglog.go

.PHONY:
build: hyperloglog.go
	go build .

.PHONY:
clean: hyperloglog
	rm -rf hyperloglog
