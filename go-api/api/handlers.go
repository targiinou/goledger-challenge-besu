package api

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"
)

type Blockchain interface {
	GetCurrentValue(ctx context.Context) (*big.Int, error)
	SetValue(ctx context.Context, newValue *big.Int) (string, error)
}

type Storage interface {
	StoreValue(ctx context.Context, value *big.Int) error
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	value, err := s.blockchain.GetCurrentValue(ctx)
	if err != nil {
		log.Printf("Erro ao buscar valor do contrato: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "O valor atual no smart contract é: %s", value.String())
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {
	valueStr := r.URL.Query().Get("value")
	if valueStr == "" {
		http.Error(w, "Parâmetro 'value' é obrigatório", http.StatusBadRequest)
		return
	}

	newValue, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		http.Error(w, "Valor inválido. Por favor, forneça um número.", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	txHash, err := s.blockchain.SetValue(ctx, newValue)
	if err != nil {
		log.Printf("Erro ao definir o valor no contrato: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Valor alterado com sucesso! Hash da transação: %s", txHash)
}

func (s *Server) syncHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	value, err := s.blockchain.GetCurrentValue(ctx)
	if err != nil {
		log.Printf("Sync - Erro ao buscar valor do contrato: %v", err)
		http.Error(w, "Erro ao buscar valor da blockchain", http.StatusInternalServerError)
		return
	}

	if err := s.storage.StoreValue(ctx, value); err != nil {
		log.Printf("Sync - Erro ao salvar valor no banco de dados: %v", err)
		http.Error(w, "Erro ao salvar valor no banco de dados", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Valor (%s) sincronizado com sucesso para o banco de dados!", value.String())
}
