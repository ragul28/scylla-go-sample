# Scylla Golang Sample

## Getting started Setup

* Initial setup of scylladb 
```sh
docker-compose up -d
docker exec -it scylla cqlsh
```

* Run with cqlsh
```sql
CREATE KEYSPACE IF NOT EXISTS sample WITH REPLICATION = {'class': 'SimpleStrategy','replication_factor':1};
describe keyspaces;
use sample;
```

* Create table inside the DB
```sql
CREATE TABLE user_data (
    id int,
    first_name text,
    last_name text,
    address text,
    PRIMARY KEY((first_name, last_name)));
```