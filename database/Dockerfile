FROM postgres:15.3-alpine3.18 AS production
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB cinema_ticket_db
COPY ./migrations/init_tables.sql docker-entrypoint-initdb.d/init_tables.sql
