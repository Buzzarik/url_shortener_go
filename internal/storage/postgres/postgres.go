package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(cnf config.ConfigPostgres) (*Postgres, error) {
	const op = "StoragePostgres.New"

	//создаем конфиг строку
	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cnf.Host,
		cnf.Port,
		cnf.Username,
		cnf.Password,
		cnf.DbName,
		cnf.Sslmode)

	db, err := sql.Open(cnf.Driver, conn_str)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(int(cnf.MaxOpenConns))
	db.SetMaxIdleConns(int(cnf.MaxIdleConns))
	db.SetConnMaxIdleTime(cnf.MaxIdleTime)

	//NOTE: проверка на подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{db: db}, nil
}


func (s *Postgres) SaveUrl(url *models.Url) error{
	const op = "Postgres.SaveUrl";

	stmt, err := s.db.Prepare(`
		INSERT INTO urls (url, alias)
		VALUES ($1, $2)
	`);

	if err != nil {
		return fmt.Errorf("(Prepare) %s, %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second);
	defer cancel();

	err = stmt.QueryRowContext(ctx, url.UrlText, url.Alias).Err();

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return fmt.Errorf("CONSTRAINT");
	}

	if err != nil {
		return fmt.Errorf("(Query) %s, %w", op, err);
	}

	return nil;
}

func (s *Postgres) GetUrl(alias string) (*models.Url, error) {
	const op = "Postgres.GetUrl";
	
	stmt, err := s.db.Prepare(`
		SELECT id, url, alias FROM urls
		WHERE alias = $1
	`);

	if err != nil {
		return nil, fmt.Errorf("(Prepare) %s, %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second);
	defer cancel();

	var url models.Url;
	err = stmt.QueryRowContext(ctx, alias).Scan(
		&url.Id,
		&url.UrlText,
		&url.Alias,
	);

	if err == sql.ErrNoRows {
		return nil, nil;
	}

	if err != nil {
		return nil, fmt.Errorf("(Query) %s, %w", op, err);
	}
	return &url, nil;
}