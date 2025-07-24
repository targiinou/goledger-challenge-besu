# GoLedger Challenge - Besu Edition

Este repositório contém a solução para o desafio GoLedger Besu. A solução consiste em uma API REST em Go que interage com uma rede blockchain Hyperledger Besu local e um banco de dados PostgreSQL para gerenciar o estado de um smart contract.

## Arquitetura da Solução

A aplicação em Go está localizada no diretório `go-api` e foi estruturada para separar as responsabilidades em pacotes distintos:

-   **`config`**: Gerencia o carregamento de configurações a partir de um arquivo `.env`, evitando que dados sensíveis (como chaves privadas) sejam expostos no código.
-   **`blockchain`**: Encapsula toda a lógica de interação com a rede Besu, incluindo a conexão com o nó, chamadas de leitura (`get`) e transações de escrita (`set`).
-   **`storage`**: Responsável pela comunicação com o banco de dados PostgreSQL. Contém a lógica para criar a tabela, salvar e buscar o valor do contrato.
-   **`api`**: Expõe a funcionalidade da aplicação através de uma API REST. Gerencia as rotas e os handlers HTTP, orquestrando as chamadas para os pacotes `blockchain` e `storage`.
-   **`main.go`**: É o ponto de entrada da aplicação. Sua única responsabilidade é inicializar os módulos de configuração, blockchain, storage e o servidor da API.

## Pré-requisitos

Antes de iniciar, certifique-se de que você tem as seguintes ferramentas instaladas e configuradas em seu ambiente (preferencialmente WSL 2 para usuários Windows):

-   [**Go**](https://go.dev/dl/) (versão 1.18+)
-   [**Docker e Docker Compose**](https://www.docker.com/products/docker-desktop/)
-   [**Node.js e NPM**](https://nodejs.org/en/download/)
-   [**Hyperledger Besu**](https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution) (com Java 21+)
-   **`jq`**: Utilitário de linha de comando para processar JSON.
    ```bash
    sudo apt-get install jq
    ```

## Como Executar a Solução

Siga os passos abaixo para configurar e executar todo o ambiente.

### 1. Iniciar a Rede Besu e Implantar o Contrato

O primeiro passo é iniciar a rede blockchain local. O script `startDev.sh` automatiza todo o processo.

```bash
# 1. Navegue até o diretório besu
cd besu

# 2. Execute o script de inicialização
# (Pode ser necessário executar 'sudo rm -rf node/besu-*/data' antes se houver resquícios de uma execução anterior)
./startDev.sh
```

Ao final da execução, o script exibirá o **endereço do smart contract implantado**. Copie este endereço, pois ele será necessário no próximo passo.

### 2. Configurar e Iniciar a API em Go

Com a rede Besu no ar, agora podemos configurar e iniciar a API.

1.  **Navegue até o diretório da API:**
    ```bash
    cd go-api
    ```

2.  **Crie o arquivo de configuração `.env`:**
    Copie o arquivo de exemplo para criar seu arquivo de configuração local.
    ```bash
    cp .env.example .env
    ```

3.  **Edite o arquivo `.env`:**
    Abra o arquivo `.env` e preencha o `CONTRACT_ADDRESS` com o endereço que você copiou no passo anterior. A chave privada (`SIGNER_PRIVATE_KEY`) já está preenchida com a chave de desenvolvimento padrão.

4.  **Inicie o Banco de Dados:**
    Use o Docker Compose para iniciar o contêiner do PostgreSQL em segundo plano.
    ```bash
    docker-compose up -d
    ```

5.  **Instale as dependências Go:**
    ```bash
    go mod tidy
    ```

6.  **Inicie a API:**
    ```bash
    go run .
    ```

O servidor estará rodando em `http://localhost:8080`.

## Endpoints da API

Você pode interagir com a API usando `curl` ou seu navegador.

-   **GET `/get`**
    Consulta o valor atual da variável no smart contract.
    ```bash
    curl http://localhost:8080/get
    ```

-   **POST `/set`**
    Altera o valor da variável no smart contract. Requer o parâmetro `value`.
    ```bash
    curl -X POST "http://localhost:8080/set?value=42"
    ```

-   **POST `/sync`**
    Busca o valor atual da blockchain e o salva/atualiza no banco de dados.
    ```bash
    curl -X POST http://localhost:8080/sync
    ```

-   **GET `/check`**
    Compara o valor na blockchain com o valor no banco de dados e retorna um JSON.
    ```bash
    curl http://localhost:8080/check
    # Exemplo de Resposta: {"are_values_equal":true,"blockchain_value":"42","database_value":"42"}
    ```

## Testando com Postman

Para facilitar os testes, uma coleção do Postman está incluída no projeto.

1.  Abra o Postman.
2.  Clique em **Import** > **File**.
3.  Navegue até o diretório `go-api/postman` e selecione o arquivo `GoLedger_Besu_Challenge.postman_collection.json`.
4.  Uma nova coleção chamada "GoLedger Besu Challenge" aparecerá, contendo as 4 requisições prontas para serem usadas.
