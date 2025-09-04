# Cinema Ticket Booking API

## About

A Cinema Ticket Booking API for case study using Go, Gin, Swaggo, and PostgreSQL. Main features for this case study are:

- User authentication and authorization with JWT
- CRUD operations for movie screenings

## Personal

Made by [Tirza Sarwono](https://www.linkedin.com/in/tirzasrwn/)

## Stack

- Go 1.25
- Gin Web Framework
- Swagger API Documentation with Swaggo
- JWT Authentication
- PostgreSQL
- Docker & Docker Compose

## Entity Relationship Diagram (ERD)

![cinema-ticket-erd](./docs/cinema-erd.png)

## Requirement

- Unix based OS (for make command)
- Docker
- Docker Compose
- Make
- Go (for development)

## Running

### Using Docker Compose (Recommended)

```sh
# Start all services
make up
# Stop all services
make down
```

### Running Specific Services

You can run individual services separately:

```sh
# Build and start specific services
make docker_db_build    # Database
make docker_db_start
make docker_db_stop

make docker_be_build    # Backend
make docker_be_start
make docker_be_stop
```

## API Routes

![swagger](./docs/swagger-ui.png)

This API documentation uses Swagger. Here are the main routes:

| Route              | Method | Description                            | Authentication |
| ------------------ | ------ | -------------------------------------- | -------------- |
| `/login`           | POST   | Authenticate user and return JWT token | Public         |
| `/screenings`      | GET    | Get all available screenings           | JWT Required   |
| `/screenings`      | POST   | Create new screening                   | JWT + Admin    |
| `/screenings/{id}` | GET    | Get specific screening details         | JWT Required   |
| `/screenings/{id}` | PUT    | Update screening information           | JWT + Admin    |
| `/screenings/{id}` | DELETE | Delete screening                       | JWT + Admin    |

Swagger API documentation can be found at [http://localhost:4000/swagger/index.html](http://localhost:4000/swagger/index.html).

Postman Collection is available in the [`./docs/`](./docs/Cinema-Ticket-API.postman_collection.json) directory.

### Feature and Route Correlations

- **User Authentication**
  - Login: `POST /login`

- **Admin Operations**
  - Manage screenings: All `/screenings` endpoints

## Service Details

### Ports

| Service    | Port | Description     |
| ---------- | ---- | --------------- |
| Backend    | 4000 | Main API server |
| PostgreSQL | 5432 | Database        |

### Default Credentials

- **API Test User**

  ```sh
  email: admin@cinema.com
  password: password
  ```

- **Database**

  ```sh
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: cinema_ticket_db
  ```

## Database Migration

Database migrations are automatically applied when starting with Docker Compose.
