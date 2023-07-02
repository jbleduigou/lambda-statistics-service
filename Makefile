.PHONY: build-ListFunctionsFunction build-SearchFunctionsFunction clean


build-SearchFunctionsFunction:
	go get -v -t -d ./...
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -o ./search-functions/bootstrap ./search-functions
	cp ./search-functions/bootstrap $(ARTIFACTS_DIR)/bootstrap

build-ListFunctionsFunction:
	go get -v -t -d ./...
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -o ./list-functions/bootstrap ./list-functions
	cp ./list-functions/bootstrap $(ARTIFACTS_DIR)/bootstrap


clean:
	rm -rf ./search-functions/bootstrap ./list-functions/bootstrap