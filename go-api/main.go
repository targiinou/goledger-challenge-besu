package main

import (
	"log"

	"github.com/seu-usuario/goledger-challenge-besu/go-api/api"
	"github.com/seu-usuario/goledger-challenge-besu/go-api/blockchain"
	"github.com/seu-usuario/goledger-challenge-besu/go-api/config"
	"github.com/seu-usuario/goledger-challenge-besu/go-api/storage"
)

func main() {
	// Carrega as configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Falha ao carregar as configurações: %v", err)
	}

	// Inicializa o cliente da blockchain com as configurações
	bcClient, err := blockchain.NewClient(cfg.BesuNodeURL, cfg.ContractAddress, cfg.SignerPrivateKey)
	if err != nil {
		log.Fatalf("Falha ao inicializar o cliente da blockchain: %v", err)
	}

	// Inicializa o storage
	dbStorage, err := storage.NewPostgresStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao inicializar o storage: %v", err)
	}

	// Cria e inicia o servidor da API
	server := api.NewServer(bcClient, dbStorage)
	server.Start()
}
