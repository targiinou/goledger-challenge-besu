package storage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(databaseURL string) (*PostgresStorage, error) {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("não foi possível pingar o banco de dados: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
	return &PostgresStorage{db: db}, nil
}

func createTable(db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_value (
		id SERIAL PRIMARY KEY,
		value TEXT NOT NULL,
		last_updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	INSERT INTO contract_value (id, value) VALUES (1, '0') ON CONFLICT (id) DO NOTHING;
	`
	_, err := db.Exec(context.Background(), query)
	return err
}

// StoreValue salva ou atualiza o valor do contrato no banco de dados.
func (s *PostgresStorage) StoreValue(ctx context.Context, value *big.Int) error {
	query := `UPDATE contract_value SET value = $1, last_updated_at = NOW() WHERE id = 1;`
	_, err := s.db.Exec(ctx, query, value.String())
	return err
}

// GetStoredValue busca o valor atual armazenado no banco de dados.
func (s *PostgresStorage) GetStoredValue(ctx context.Context) (*big.Int, error) {
	var valueStr string
	query := `SELECT value FROM contract_value WHERE id = 1;`

	err := s.db.QueryRow(ctx, query).Scan(&valueStr)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar valor no banco de dados: %w", err)
	}

	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return nil, fmt.Errorf("falha ao converter valor do banco de dados para big.Int")
	}

	return value, nil
}
