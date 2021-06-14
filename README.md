# BookstoreAPI

# Docker Instructions

1: Download source files and extract

2: Open terminal in project directory

3: Enter Following Build and Run commands

# Build terminal command:

docker build -t my-go-app .

# Run Project as well as all unit tests

docker run -p 8080:8081 my-go-app sh -c "go test"
