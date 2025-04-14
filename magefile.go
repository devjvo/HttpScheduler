//go:build mage
// +build mage

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
)

func Init() {
	Up()
	WaitForDb()
	DbMigrateUp()
	LoadFixtures()
}

func Up() error {
	return execute(
		exec.Command("docker", "compose", "-f", "docker-compose.yaml", "up", "-d"),
		"Up containers...",
	)
}

func Down() error {
	return execute(
		exec.Command("docker", "compose", "-f", "docker-compose.yaml", "down"),
		"Down containers...",
	)
}

func Ps() error {
	return execute(
		exec.Command("docker", "compose", "-f", "docker-compose.yaml", "ps"),
		"",
	)
}

func WaitForDb() error {
	timeoutAt := time.Now().Add(10000000000) // 10s

	for {
		db, err := sql.Open("postgres", "dbname=http_scheduler user=dev password=dev port=2432 sslmode=disable")

		if err != nil {
			return fmt.Errorf(err.Error())
		}

		context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err = db.PingContext(context)
		cancel() // prevent time to keep running in background

		if err == nil {
			fmt.Println("Database is up")
			return nil
		}

		db.Close()

		if time.Now().After(timeoutAt) {
			return fmt.Errorf("Timeout! Database is still unavailable")
		}

		fmt.Println("Database not available. Let's try again...")
		time.Sleep((1 * time.Second))
	}

	return nil
}

func DbConnect() error {
	return execute(
		exec.Command("psql", "-h", "127.0.0.1", "-p", "2432", "-U", "dev", "-d", "http_scheduler"),
		"Connecting to the PostgreSQL database...",
	)
}

func MigrateCreate(tableName string) error {
	return execute(
		exec.Command("migrate", "create", "-ext", "sql", "-dir", "db/migrations", "-seq", tableName),
		"",
	)
}

func DbMigrateUp() error {
	return databaseMigrate("up")
}

func DbMigrateDown() error {
	return databaseMigrate("down")
}

func databaseMigrate(direction string) error {
	return execute(
		exec.Command("migrate", "-database", "postgres://dev:dev@127.0.0.1:2432/http_scheduler?sslmode=disable", "-path", "db/migrations", direction),
		fmt.Sprintf("Migration %s", direction),
	)
}

func LoadFixtures() error {
	fmt.Println("*****")
	fmt.Println("Loading fixtures...")
	var err error
	var db *sql.DB
	var fixtures *testfixtures.Loader

	db, err = sql.Open("postgres", "dbname=http_scheduler user=dev password=dev port=2432 sslmode=disable")

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	fixtures, err = testfixtures.New(
		testfixtures.DangerousSkipTestDatabaseCheck(), // https://github.com/go-testfixtures/testfixtures?tab=readme-ov-file#security-check
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("db/fixtures"),
	)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = fixtures.Load()

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fmt.Println("Fixtures loaded.")

	return nil
}

func Shell() error {

	cmd := exec.Command("docker", "compose", "exec", "api", "sh")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error connecting to the container of the API: %v", err)
	}

	return nil
}

func execute(commandToExec *exec.Cmd, message string) error {
	if message != "" {
		fmt.Println("*****")
		fmt.Println(message)
	}

	commandToExec.Stdout = os.Stdout
	commandToExec.Stderr = os.Stderr
	commandToExec.Stdin = os.Stdin

	err := commandToExec.Run()

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
