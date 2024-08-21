package psg

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sync"
	"taskTest/internal/config"
	"taskTest/internal/entity"
)

type PostgresStorage struct {
	db *sql.DB
	mu *sync.RWMutex
}

func NewPostgresStorage(conf *config.Config) (*PostgresStorage, error) {

	login := conf.DB_LOGIN
	password := conf.DB_PASS
	host := conf.DB_HOST
	port := conf.DB_PORT
	dbname := conf.DB_NAME

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", login, password, host, port, dbname)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Создание таблицы
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Create(ctx context.Context, user *entity.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStorage) Read(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE id=$1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStorage) Update(ctx context.Context, user entity.User) error {
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	_, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	return err
}

func (s *PostgresStorage) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
