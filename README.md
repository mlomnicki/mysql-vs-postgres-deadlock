Synopsis
========

It shows how a innocent transaction deadlocks MySQL but works perfectly well on Postgres.

See https://bugs.mysql.com/bug.php?id=48652

    $ go run deadlock.go mysql "user:password@tcp(127.0.0.1)/testDb"

    Error 1213: Deadlock found when trying to get lock; try restarting transaction
    Error 1213: Deadlock found when trying to get lock; try restarting transaction
    Error 1213: Deadlock found when trying to get lock; try restarting transaction
    Error 1213: Deadlock found when trying to get lock; try restarting transaction
    Error 1213: Deadlock found when trying to get lock; try restarting transaction
    Error 1213: Deadlock found when trying to get lock; try restarting transaction

    $ go run deadlock.go postgresql "postgres://user:password@localhost/dbname"

    No deadlocks!

Connection strings
==================

- Postgres: see https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
- MySql: see https://github.com/go-sql-driver/mysql#dsn-data-source-name
