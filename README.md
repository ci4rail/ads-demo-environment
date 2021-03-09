# ads-demo-environment
This repository contains the components of ads-demo-environment including their deplyoment to aks.

## Test environment
### TimescaleDB Server

Start server with
```
docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb:2.1.0-pg13
```

### TimescaleDb Go Client
Run local go program to
* Extend the database with TimescaleDB (if not happed before)
* Create table if not exists
* Insert columns if not exist
* Convert the table created into a hypertable (if not exists)
* Insert sample data into database

Precoditions:
* local docker container stared as noted in [TimescaleDB Server](#TimescaleDB-Server)

Execute it
```
go run main.go
```

### TimescaleDB Docker Client
Run interactive client in docker container and connect to PostgreSQL, using a superuser named 'postgres':
```
docker exec -it timescaledb psql -U postgres
```

Show the last 100 entries in the table:
```
SELECT * FROM adsdata ORDER BY time DESC LIMIT 100;
```

### Grafana
Run as docker container:
```
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```
