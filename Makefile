.PHONY: test

test:
	go test -v -race -short -count=1 -coverprofile=.cover ./... 

cover:
	go tool cover -html=.cover