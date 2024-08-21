package psg

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"taskTest/internal/config"
	"taskTest/internal/entity"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func InitStorage(conf *config.Config) (*PostgresStorage, error) {

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

	// Чтение schema.sql
	schema, err := os.ReadFile("migrations/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения schema.sql: %v", err)
	}

	// Выполнение SQL-запросов
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения schema.sql: %v", err)
	}

	return NewPostgresStorage(db), nil
}

func (s *PostgresStorage) CreateUser(user *entity.User) error {

	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := s.db.Query(query, user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("Error CreateUser: %v", err)
	}
	return nil
}

func (s *PostgresStorage) GetUserByID(id int) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE id=$1`
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStorage) UpdateUser(user *entity.User) error {
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	_, err := s.db.Exec(query, user.Name, user.Email, user.ID)
	return err
}

func (s *PostgresStorage) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := s.db.Exec(query, id)
	return err
}
