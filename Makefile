
build:
	go generate ./...
	go build -o sliverStager ./cmd/sliverStager