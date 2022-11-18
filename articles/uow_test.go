package articles_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/antunesleo/rest-api-go/articles"
	"github.com/go-pg/pg/v10"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	migrate "github.com/rubenv/sql-migrate"
)

func applyMigrations(t *testing.T, host, port, user, password, DB string) error {
	t.Helper()
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		DB,
	)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	migrations := &migrate.FileMigrationSource{
		Dir: "../migrations",
	}
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}

func setupTestDatabase(t *testing.T, DB, password, user string) (testcontainers.Container, *pg.DB, error) {
	t.Helper()
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       DB,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_USER":     user,
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	if err != nil {
		fmt.Println("error creating container", err)
		return nil, nil, err
	}

	// 3.2 Create db connection string and connect
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: DB,
		Addr:     host + ":" + port.Port(),
	})

	err = applyMigrations(t, host, port.Port(), user, password, DB)
	if err != nil {
		fmt.Println("Error applying migration", err)
	}

	return dbContainer, db, nil
}

func TestUnitOfWork(t *testing.T) {
	ctx := context.TODO()

	t.Run("WithinTx should commit transaction", func(t *testing.T) {
		postgresDB := "testdb"
		postgresPassword := "postgres"
		postgresUser := "postgres"
		dbContainer, db, err := setupTestDatabase(t, postgresDB, postgresPassword, postgresUser)
		if err != nil {
			t.Error(err)
		}
		defer dbContainer.Terminate(context.Background())

		uow := articles.NewGoPgUow(db)

		categoryName := "category"

		createStuffInTx := func(uows articles.UowStore) error {
			category := articles.Category{Name: categoryName}
			return uows.GetCategoryRepository().Add(&category)
		}
		err = uow.WithinTx(ctx, createStuffInTx)
		assert.NoError(t, err)

		categoryDB := &articles.CategoryDB{}
		err = db.Model(categoryDB).
			Where("name = ?", categoryName).
			Select()

		assert.NoError(t, err)
	})

	t.Run("WithinTx should rollback transaction on error", func(t *testing.T) {
		postgresDB := "testdb"
		postgresPassword := "postgres"
		postgresUser := "postgres"
		dbContainer, db, err := setupTestDatabase(t, postgresDB, postgresPassword, postgresUser)
		if err != nil {
			t.Error(err)
		}
		defer dbContainer.Terminate(context.Background())

		uow := articles.NewGoPgUow(db)

		categoryName := "category"

		rollbackStuffInTx := func(uows articles.UowStore) error {
			category := articles.Category{Name: categoryName}
			uows.GetCategoryRepository().Add(&category)
			return errors.New("some random error")
		}
		err = uow.WithinTx(ctx, rollbackStuffInTx)
		assert.Error(t, err)

		categoryDB := &articles.CategoryDB{}
		err = db.Model(categoryDB).
			Where("name = ?", categoryName).
			Select()

		assert.Error(t, err)
	})
}
