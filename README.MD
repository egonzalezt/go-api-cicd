# Go APi CI/CD Eafit 2023-2

## Table of Contents
- [Getting Started](#getting-started)
  - [Installing Dependencies](#installing-dependencies)
  - [Running the Service](#running-the-service)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Running with Docker](#running-with-docker)
- [Environment Variables](#environment-variables)
- [License](#license)

## Getting Started

### Installing Dependencies

To set up the required Go dependencies, run the following command in your project's root directory:

```bash
go mod tidy
```

## Running the Service

You can start the API service by running the following command:

```bash
go run main.go
```

# API Endpoints

The API provides the following endpoints for managing notes:

## Get All Notes and Create a New Note

* GET All Notes and POST New Note: http://localhost:3001/api/notes
## Update and Get Note by ID
* PUT Note by ID and GET Note by ID: http://localhost:3001/api/notes/{id}
## Home Page
* Home Page: http://localhost:3001/

# Testing

To run tests and generate a coverage report, you can use the following commands:

Generate a coverage report:
```bash
go test -coverprofile=coverage.txt -covermode count ./
```

Generate an XML report:

```bash
/home/$USER/go/bin/gocover-cobertura < coverage.txt > coverage.xml
```

# Running with Docker

This code supports Docker, making it easy to containerize your application. Follow these steps to run the application using Docker:

## Steps

1. Build and start the containers:
```bash
docker-compose up --build -d
```
2. Open your web browser and access the application at localhost:5000.
If you wish to remove networks and volumes associated with this project, you can use the following command:

```bash
docker-compose down -v
```

# Environment Variables

The application reads the following environment variable:

`PORT`: You can set the port on which the service runs. If not set, it defaults to `5000`.

# License

This project is licensed under the MIT License - see the [LICENSE](./README.md) file for details.