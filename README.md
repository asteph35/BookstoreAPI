# BookstoreAPI

# Docker Instructions

# Build terminal command:

docker build -t my-go-app .

# Run Project as well as all unit tests

docker run -p 8080:8081 my-go-app sh -c "go test"
