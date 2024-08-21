package usecase

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"taskTest/internal/entity"
	"taskTest/internal/storage/psg"
	"testing"
)

func setupPostgresContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// PostgreSQL контейнер
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_pas",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("postgres://test_user:test_pas@%s:%s/test_db?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	schema, err := os.ReadFile("../../migrations/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Выполнение SQL-запросов
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}

// добавление записи, добавление дубликатов уникальных полей
func TestCreateUser(t *testing.T) {
	db, teardown := setupPostgresContainer(t)
	defer teardown()

	repo := psg.NewPostgresStorage(db)
	uc := NewUseCase(repo)

	//проверяем количество записей
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, count)

	// Тест CreateUser
	user := &entity.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	err = uc.CreateUser(user)
	assert.NoError(t, err)

	//проверяем количество записей
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, count)

	//добавляем вторую запись
	user2 := &entity.User{ID: 56, Name: "John Uik", Email: "johnuik@ex.ru"}
	err = uc.CreateUser(user2)
	assert.NoError(t, err)

	//проверяем количество записей
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, count)

	//добавляем запись с повторяющимся ID
	user3 := &entity.User{ID: 1, Name: "John Uik", Email: "johnu33ik@ex.ru"}
	err = uc.CreateUser(user3)
	assert.Error(t, err)
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, count)

	//добавляем запись с повторяющимся email
	user4 := &entity.User{ID: 344, Name: "John Uik", Email: "johnuik@ex.ru"}
	err = uc.CreateUser(user4)
	assert.Error(t, err)
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, count)

}

func TestGetUserByID(t *testing.T) {
	db, teardown := setupPostgresContainer(t)
	defer teardown()

	repo := psg.NewPostgresStorage(db)
	uc := NewUseCase(repo)

	//наполнение бд тестовыми данными
	user := &entity.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	err := uc.CreateUser(user)
	assert.NoError(t, err)
	user2 := &entity.User{ID: 2, Name: "Ser Goreliy", Email: "serga@mail.co"}
	err = uc.CreateUser(user2)
	assert.NoError(t, err)

	//обычный запрос 1
	savedUser, err := uc.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, savedUser)
	assert.Equal(t, user.Name, savedUser.Name)
	assert.Equal(t, user.Email, savedUser.Email)

	//обычный запрос 2
	savedUser, err = uc.GetUserByID(user2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, savedUser)
	assert.Equal(t, user2.Name, savedUser.Name)
	assert.Equal(t, user2.Email, savedUser.Email)

	//запрос с несуществующим ID
	savedUser, err = uc.GetUserByID(15)
	assert.Error(t, err)
	assert.Nil(t, savedUser)

}

func TestUpdateUser(t *testing.T) {
	db, teardown := setupPostgresContainer(t)
	defer teardown()

	repo := psg.NewPostgresStorage(db)
	uc := NewUseCase(repo)

	//наполнение бд тестовыми данными
	user := &entity.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	err := uc.CreateUser(user)
	assert.NoError(t, err)
	user2 := &entity.User{ID: 2, Name: "Ser Goreliy", Email: "serga@mail.co"}
	err = uc.CreateUser(user2)
	assert.NoError(t, err)

	// Тест 1
	user.Name = "John Smith"
	err = uc.UpdateUser(user)
	assert.NoError(t, err)
	updatedUser, err := uc.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, updatedUser.Name)

	// Тест 2
	user2.Name = "Guffy"
	err = uc.UpdateUser(user2)
	assert.NoError(t, err)
	updatedUser, err = uc.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, updatedUser.Name)

}

func TestDeleteUser(t *testing.T) {
	db, teardown := setupPostgresContainer(t)
	defer teardown()

	repo := psg.NewPostgresStorage(db)
	uc := NewUseCase(repo)

	//наполнение бд тестовыми данными
	user := &entity.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	err := uc.CreateUser(user)
	assert.NoError(t, err)
	user2 := &entity.User{ID: 2, Name: "Ser Goreliy", Email: "serga@mail.co"}
	err = uc.CreateUser(user2)
	assert.NoError(t, err)

	//проверяем количество записей
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, count)

	// Тест DeleteUser 1
	err = uc.DeleteUser(user.ID)
	assert.NoError(t, err)
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, count)
	deletedUser, err := uc.GetUserByID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)

	// Тест DeleteUser 2
	err = uc.DeleteUser(user2.ID)
	assert.NoError(t, err)
	err = db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, count)
	deletedUser, err = uc.GetUserByID(user2.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}
