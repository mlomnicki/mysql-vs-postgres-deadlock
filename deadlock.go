package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"sync"
)

const CONCURRENCY_LEVEL = 30

const MySQL_SCHEMA = `
    DROP TABLE IF EXISTS comments;
    DROP TABLE IF EXISTS posts;
    CREATE TABLE posts (
      id             INT(11) NOT NULL AUTO_INCREMENT,
      comments_count INT(11) NOT NULL,
      PRIMARY KEY (id)
    );
    CREATE TABLE comments (
      id      INT(11) NOT NULL AUTO_INCREMENT,
      post_id INT(11) NOT NULL,
      PRIMARY KEY (id),
      FOREIGN KEY (post_id) REFERENCES posts(id)
    );
    INSERT INTO posts (id, comments_count) VALUES (1, 0);
`

const POSTGRES_SCHEMA = `
  DROP TABLE IF EXISTS comments;
  DROP TABLE IF EXISTS posts;
  CREATE TABLE posts (
    id             serial PRIMARY KEY,
    comments_count integer
  );
  CREATE TABLE comments (
    id      serial PRIMARY KEY,
    post_id integer REFERENCES posts
  );
  INSERT INTO posts (id, comments_count) VALUES (1, 0);
`

const DEADLOCKING_STATEMENT = `
	BEGIN;
  INSERT INTO comments (post_id) VALUES (1);
  UPDATE posts SET comments_count = comments_count + 1 WHERE id = 1;
	COMMIT;
`

func simulateDeadlock(fn func()) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(CONCURRENCY_LEVEL)
	for i := 0; i < CONCURRENCY_LEVEL; i++ {
		go func() {
			<-start
			fn()
			wg.Done()
		}()
	}
	close(start)
	wg.Wait()
}

func mysql(connectionString string) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mysqlExec(db, MySQL_SCHEMA)
	simulateDeadlock(func() { mysqlExec(db, DEADLOCKING_STATEMENT) })
}

func postgres(connectionString string) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Query(POSTGRES_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
	simulateDeadlock(func() { postgresExec(db, DEADLOCKING_STATEMENT) })
}

func mysqlExec(db *sql.DB, statement string) {
	for _, s := range strings.Split(statement, ";") {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) > 0 {
			_, err := db.Exec(trimmed)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func postgresExec(db *sql.DB, statement string) {
	_, err := db.Query(statement)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Printf("Usage: %s mysql|postgres\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	dbEngine, connectionString := os.Args[1], os.Args[2]
	switch dbEngine {
	case "mysql":
		mysql(connectionString)
	case "postgres":
		postgres(connectionString)
	default:
		usage()
	}
}
