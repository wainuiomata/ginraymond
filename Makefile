.PHONY : build
build :
	@go build

.PHONY : update
update :
	@go get -u
	@go mod tidy

.PHONY : lint
lint :
	@golangci-lint run

.PHONY : test
test :
	@go test -coverprofile=coverage.txt -covermode count
	@go tool cover -func=coverage.txt
	@go tool cover -html=coverage.txt -o=coverage-report.html
