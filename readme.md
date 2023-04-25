
## go-fetcher

fetcher tool, sql-data to PostgresSQL

### installation

#### PostgresSQL

```
sudo apt-get install wget ca-certificates
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get install postgresql

psql -V

sudo systemctl enable postgresql
sudo systemctl start postgresql
sudo systemctl status postgresql
```


#### Config

/etc/postgresql/12/main/postgresql.conf
```
# - Connection Settings -
listen_addresses='*'
#listen_addresses = 'localhost'
```

/etc/postgresql/12/main/pg_hba.conf
```
// use password login
    # TYPE  DATABASE        USER            ADDRESS                 METHOD
    host    all             all             0.0.0.0/0               md5

// use no password
    # TYPE  DATABASE        USER            ADDRESS                 METHOD
    host    all             all             0.0.0.0/0               trust
```

Restart postgresql
```
sudo systemctl restart postgresql
```


#### db create
```
// sql check
sudo -u postgres psql fabricexplorer
\l
\d

// Option: eunc database and tables config
// Create user and database
CREATE USER eunoiaAd WITH PASSWORD 'Trusme123#@!';
CREATE DATABASE eunoiaDB OWNER eunoiaAd;
GRANT ALL PRIVILEGES ON DATABASE eunoiaDB TO eunoiaAd;

// Create tables: blocks and transactions
\i eunc.sql
```


### Usage

```

// fetch
$ go build fetcher.go && ./fetcher -T 32 -s eunc -d localhost:5432/eunoiaDB | tee $(date "+%Y%m%d_%H%M%S.log")

// postgresql config
- database: eunoiaDB
- user: eunoiaAd
- password: Trusme123#@!
- port 5432

// postgresql connect by terminal
$ psql postgres://eunoiaAd:Trusme123#@!@localhost:5432/eunoiaDB
$ psql -U eunoiaAd -W -h localhost -d eunoiaDB -p 5432

// sql server explorer
- http://localhost/phppgadmin/index.php
```