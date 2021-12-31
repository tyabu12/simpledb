TARGET:=simpledb.exe
OBJ:=$(shell find . -type f -name "*.go")

$(TARGET): $(OBJ)
	go build -o $@

.PHONY: test
test:
	go test -v ./...

.PHONY: ci-test
ci-test:
	go test -v -coverprofile=cover.out ./...

.PHONY: clean
clean:
	go clean
	rm -rf $(TARGET) cover.out