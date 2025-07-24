package api

import (
	"fmt"
	"log"
	"net/http"
)

// Server contém as dependências do servidor da API
type Server struct {
	blockchain Blockchain
	storage    Storage
}

// NewServer cria uma nova instância do servidor
func NewServer(bc Blockchain, st Storage) *Server {
	return &Server{
		blockchain: bc,
		storage:    st,
	}
}

// Start inicia o servidor HTTP e registra as rotas
func (s *Server) Start() {
	http.HandleFunc("/get", s.getHandler)
	http.HandleFunc("/set", s.setHandler)
	http.HandleFunc("/sync", s.syncHandler)

	fmt.Println("Servidor da API iniciado na porta 8080.")
	fmt.Println("Para consultar, acesse: http://localhost:8080/get")
	fmt.Println("Para alterar, acesse: http://localhost:8080/set?value=SEU_NUMERO")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
