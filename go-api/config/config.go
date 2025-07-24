package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	BesuNodeURL      string
	ContractAddress  string
	SignerPrivateKey string
	DatabaseURL      string
}

// Load carrega as configurações do arquivo .env
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		BesuNodeURL:      os.Getenv("BESU_NODE_URL"),
		ContractAddress:  os.Getenv("CONTRACT_ADDRESS"),
		SignerPrivateKey: os.Getenv("SIGNER_PRIVATE_KEY"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
	}

	return cfg, nil
}
