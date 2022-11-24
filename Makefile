lambda:
	GOARCH=amd64 GOOS=linux go build -o ./.build/main ./cmd/dgraph/dgraph.go
	zip -jrm ./.build/main.zip ./.build/main