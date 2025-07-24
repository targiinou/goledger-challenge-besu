package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const abiPath = "./abi/SimpleStorage.json"

// Client encapsula a lógica de interação com a blockchain
type Client struct {
	ethClient       *ethclient.Client
	contract        *bind.BoundContract
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
}

// NewClient cria e inicializa um novo cliente para interagir com a blockchain
func NewClient(nodeURL, contractHex, privateKeyHex string) (*Client, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com o nó Besu: %w", err)
	}

	address := common.HexToAddress(contractHex)

	abiFile, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o arquivo ABI: %w", err)
	}

	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	if err := json.Unmarshal(abiFile, &artifact); err != nil {
		return nil, fmt.Errorf("falha ao fazer o unmarshal do JSON do artefato: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(artifact.ABI)))
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer o parse do ABI do contrato: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar a chave privada: %w", err)
	}

	contract := bind.NewBoundContract(address, parsedABI, client, client, client)

	return &Client{
		ethClient:       client,
		contract:        contract,
		contractAddress: address,
		privateKey:      privateKey,
	}, nil
}

// GetCurrentValue busca o valor atual no smart contract
func (c *Client) GetCurrentValue(ctx context.Context) (*big.Int, error) {
	var results []interface{}
	callOpts := &bind.CallOpts{Context: ctx}

	err := c.contract.Call(callOpts, &results, "get")
	if err != nil {
		return nil, fmt.Errorf("falha ao chamar a função 'get' do contrato: %w", err)
	}

	value, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("não foi possível converter o resultado para *big.Int")
	}
	return value, nil
}

// SetValue envia uma transação para alterar o valor no smart contract
func (c *Client) SetValue(ctx context.Context, newValue *big.Int) (string, error) {
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("falha ao obter o chainID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("falha ao criar o transactor: %w", err)
	}

	tx, err := c.contract.Transact(auth, "set", newValue)
	if err != nil {
		return "", fmt.Errorf("falha ao criar a transação 'set': %w", err)
	}

	slog.Info("Transação enviada, aguardando mineração...", "tx", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return "", fmt.Errorf("falha ao aguardar a mineração da transação: %w", err)
	}

	slog.Info("Transação minerada com sucesso!", "receipt_status", receipt.Status)
	return tx.Hash().Hex(), nil
}
