# Backend Go for development.
start: swagger
	go run .
swagger:
	swag fmt
	swag init -g ./main.go --output ./docs/ --parseDependency

# Docker network.
docker_net_start:
	- docker network create --driver bridge cinema-ticket-api
docker_net_stop:
	- docker network remove cinema-ticket-api

# Docker backend Go.
docker_be_build:
	docker build . -t cinema-ticket-api-backend:v0.0.1
docker_be_start: docker_net_start
	docker run --network=cinema-ticket-api --name cinema-ticket-api-backend -p 4000:4000 -d cinema-ticket-api-backend:v0.0.1
docker_be_stop:
	- docker stop cinema-ticket-api-backend
	- docker rm cinema-ticket-api-backend

# Docker database PostgreSLQ.
docker_db_build:
	docker build -f ./database/Dockerfile ./database/ -t cinema-ticket-api-postgres:v0.0.1
docker_db_start: docker_net_start
	docker run --hostname cinema-ticket-api-postgres --network=cinema-ticket-api --restart unless-stopped --name cinema-ticket-api-postgres -p 5432:5432 -d cinema-ticket-api-postgres:v0.0.1
	# Use this command if you want to save the db data.
	# docker run --hostname cinema-ticket-api-postgres --network=cinema-ticket-api --restart unless-stopped --name cinema-ticket-api-postgres -p 5432:5432 -v ./database/data/.postgres/:/var/lib/postgresql/data -d cinema-ticket-api-postgres:v0.0.1
docker_db_stop:
	- docker stop cinema-ticket-api-postgres
	- docker rm cinema-ticket-api-postgres

docker_stop: docker_be_stop docker_db_stop docker_net_stop

docker_start: docker_net_start docker_db_build docker_db_start docker_be_build docker_be_start

# Docker compose backend Go and database PostgreSQL.
down:
	docker compose down
up:
	docker compose up --build -d

test:
	@rm -f coverage.out coverage.html
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o ./coverage.html
	@open "file://$(PWD)/coverage.html"
