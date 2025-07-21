package integration_test

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"testing"

// 	"github.com/ilyapiatykh/itk/config"
// 	"github.com/ory/dockertest/v3"
// )

// var db  *sql.DB
// var

// func TestMain(m *testing.M) {
// 	var err error
// 	cfg, err = config.NewCfg()
// 	if err != nil {
// 		log.Fatalf("Could not parse cfg: %s", err)
// 	}
// 	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not construct pool: %s", err)
// 	}

// 	// uses pool to try to connect to Docker
// 	err = pool.Client.Ping()
// 	if err != nil {
// 		log.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	// pulls an image, creates a container based on it and runs it
// 	resource, err := pool.Run("postgres", "latest", []string{"MYSQL_ROOT_PASSWORD=secret"})
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// 	if err := pool.Retry(func() error {
// 		var err error
// 		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to database: %s", err)
// 	}

// 	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
// 	defer func() {
// 		if err := pool.Purge(resource); err != nil {
// 			log.Fatalf("Could not purge resource: %s", err)
// 		}

// 	}()

// 	m.Run()
// }

// func TestSomething(t *testing.T) {
// 	// db.Query()
// }
