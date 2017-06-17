build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go
	docker build -t memolab/say .
	rm -f app