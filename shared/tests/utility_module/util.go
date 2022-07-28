package utility_module

import (
	"context"
	"fmt"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pokt-network/pocket/persistence"
	"github.com/pokt-network/pocket/shared/modules"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

const (
	user             = "postgres"
	password         = "secret"
	db               = "postgres"
	sql_schema       = "test_schema"
	dialect          = "postgres"
	connStringFormat = "postgres://%s:%s@%s/%s?sslmode=disable"
)

func init() {
	PersistenceModule = modules.PersistenceModule(nil) // TODO (team) make thread safe
	PostgresDB = new(persistence.PostgresDB)
}

// TODO (team) cleanup and simplify

var PostgresDB *persistence.PostgresDB
var PersistenceModule modules.PersistenceModule
var databaseUrl string

func SetupPostgresDocker(_ *testing.M) (*dockertest.Pool, *dockertest.Resource) {
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}
	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("***Make sure your docker daemon is running!!*** Could not start resource: %s\n", err.Error())
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl = fmt.Sprintf(connStringFormat, user, password, hostAndPort, db)

	log.Println("Connecting to database on url: ", databaseUrl)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("exit signal %d received\n", sig)
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("could not purge resource: %s", err)
			}
		}
	}()

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		conn, err := persistence.ConnectAndInitializeDatabase(databaseUrl, sql_schema)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		PostgresDB.Tx, err = conn.Begin(context.TODO())
		if err != nil {
			log.Println(err.Error())
			return err
		}
		PersistenceModule = persistence.NewPersistenceModule(databaseUrl, sql_schema, conn, nil)
		return nil
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}
	return pool, resource
}

func CleanupPostgresDocker(_ *testing.M, pool *dockertest.Pool, resource *dockertest.Resource) {
	defer func() {
		ctx, _ := PostgresDB.GetContext()
		PostgresDB.Tx.Rollback(ctx)
		PostgresDB.Tx = nil
		ctx.Done()
	}()
	// You can't defer this because `os.Exit`` doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(0)
}
