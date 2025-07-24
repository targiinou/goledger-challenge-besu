package main

import (
	"log"

	"github.com/seu-usuario/goledger-challenge-besu/go-api/api"
	"github.com/seu-usuario/goledger-challenge-besu/go-api/blockchain"
	"github.com/seu-usuario/goledger-challenge-besu/go-api/config"
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

	// Cria e inicia o servidor da API
	server := api.NewServer(bcClient)
	server.Start()
}
