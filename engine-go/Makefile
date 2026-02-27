test:
	go test -v ./...

cover:
	echo "Generating cover.out and cover.html"
	go test ./...  -cover cover.out
	go tool cover -html cover.out -o cover.html
