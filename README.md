# 1. Data/Request Flow
The sequence diagram illustrates these three core features:
1. Order placement
2. Driver matching
3. Live tracking

![alt text](data-request-sequence-diagram.png)

# 2. Core DB and Backend Groundwork
1. Docker Compose as the setup. Then, Postgres + PostGIS and Redis running in containers in a `docker-compose.yaml` file.
- so that we don't have to run PostGIS and Redis manually.
```

```

2. 