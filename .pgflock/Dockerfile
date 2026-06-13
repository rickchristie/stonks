# Generated-style pgflock image for local test databases.
FROM postgres:15

ENV POSTGRES_DB=postgres
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=pgflock

COPY ./postgresql.conf /etc/postgresql/postgresql.conf
ENV POSTGRES_CONFIG_FILE=/etc/postgresql/postgresql.conf

COPY ./init.sh /docker-entrypoint-initdb.d/
