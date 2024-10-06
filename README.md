## Overview

This repository contains a simple RESTful API task that facilitates the population and management of an employee search directory.

## Folder structure
```
omnihr-coding-test/
├── cmd
│  └── server
│     └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── pkg
│  ├── api
│  │  ├── employees.go
│  │  ├── employees_test.go
│  │  ├── router.go
│  │  └── user.go
│  ├── auth
│  │  ├── auth.go
│  │  ├── auth_test.go
|  |  └── utils.go   
│  ├── cache
│  │  ├── cache.go
│  │  ├── cache_mock.go
│  │  └── cache_test.go
│  ├── config
│  │  ├── config.go
│  │  ├── config.yaml
│  ├── database
│  │  ├── db.go
│  │  ├── db_mock.go
│  │  └── db_test.go
│  ├── middleware
│  │  ├── api_key.go
│  │  ├── jwt.go
│  │  ├── cors.go
│  │  ├── rate_limit.go
│  └── models
│     ├── config.go
│     ├── employee.go
│     └── user.go
└── README.md
```


## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Installation

1. Build and run the Docker containers
```bash
make build && make up
```
2. Run the services
```bash
make run
```

## Usage
### Endpoints
```bash
- GET /api/v1/employees: Get all employees.
- POST /api/v1/login: Login, auth using x-api-key.
- POST /api/v1/register: Register a new user, auth using x-api-key.
```
### Authentication
To use authenticated routes, you must include the Authorization header with the JWT token which generated when login.

```bash
curl --location 'http://localhost:8001/api/v1/employees?department=Engineering&position=DevOps%20Engineer&location=Singapore&status=Not%20Started' \
--header 'Authorization: Bearer xxx'
```


## Checklist
- [x] The service must be containerized
- [x] User authentication based on company id to avoid a data leak in the API such ( information of other organisation’s  users, extra attributes of user that is not displayed on the UI etc ) 
- [x] The API must be unit tested 
- [x] The API information must be sharable in an OPEN API format 
- [x] No external library is allowed for rate-limitting
- [x] Dynamic Columns (config stored in config.yaml)
- [x] Implemented all filter options

## Test Result
