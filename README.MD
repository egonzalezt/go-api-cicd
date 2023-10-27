# Go APi CI/CD Eafit 2023-2

# Installing dependencies
```bash
go mod tidy
```
# For running the service

```bash
go run main.go
```

# you can make a request to

## GET ALL, POST

[http://localhost:3001/api/notes](http://localhost:3001/api/notes)

## PUT, GET BY ID
[http://localhost:3001/api/notes/{id}](http://localhost:3001/api/notes/{id})

## Home page
[http://localhost:3001/](http://localhost:3001/)


# For testing 
Now generate report
```bash
go test -coverprofile=coverage.txt -covermode count ./
```
Now generate xml
```bash
/home/$USER/go/bin/gocover-cobertura < coverage.txt > coverage.xml
```

# Run docker

This code supports docker, you can run it using `docker compose`:

## Steps

Run `docker compose up --build -d` then open your browser on localhost:5000 and you will see your application

If you also wish to remove networks and volumes associated with this project, you can add the `-v` option to the command:

 ```
 docker-compose down -v
 ```
