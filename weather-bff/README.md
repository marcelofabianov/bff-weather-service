# Weather BFF Service (Serviço A)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Este serviço atua como um **Backend-for-Frontend (BFF)** / **API Gateway** para o ecossistema Weather API. Sua principal responsabilidade é prover um contrato de API público e estável, validar as requisições de entrada e orquestrar as chamadas para o `weather-service` (Serviço B), que contém a lógica de negócio principal.

## ✨ Features e Stack Tecnológica

* **Linguagem**: Go 1.25+
* **Arquitetura**: Hexagonal (Ports & Adapters), Clean Architecture.
* **Padrão de API**: BFF (Backend-for-Frontend).
* **Injeção de Dependência**: `uber-go/fx`.
* **Roteamento HTTP**: `chi`.
* **Configuração**: `viper`.
* **Logging Estruturado**: `slog`.
* **Validação**: `go-playground/validator/v10` integrado com uma package de erro customizada (`fault`).
* **Containerização**: Docker e Docker Compose.

## 🚀 Como Executar

Existem duas formas principais de executar este serviço: como parte do sistema completo (recomendado) ou de forma isolada para desenvolvimento.

### Executando o Sistema Completo (Recomendado)

A forma ideal de rodar o BFF é junto com suas dependências. Para isso, utilize o `docker-compose.yml` localizado no diretório raiz (`bff-weather-service`).

1.  A partir do **diretório raiz**, configure os arquivos `.env` para cada serviço, especialmente o `weather-service/.env` com a chave da WeatherAPI.
2.  Execute o comando:
    ```sh
    docker compose up --build
    ```
3.  O BFF estará disponível para requisições em `http://localhost:8081`.

### Executando Apenas o BFF (Para Desenvolvimento)

Para desenvolver e testar o BFF de forma isolada, você pode executá-lo diretamente.

**Pré-requisito:** É necessário ter uma instância do `weather-service` rodando e acessível, pois o BFF depende dela.

1.  Navegue até o diretório `weather-bff`.
2.  Crie seu arquivo `.env`: `cp .env.example .env`.
3.  **Importante:** Se o `weather-service` estiver rodando localmente na porta 8080, ajuste a variável `APP_CLIENTS_WEATHERSERVICE_URL` no `.env` para `http://localhost:8080`.
4.  Exporte as variáveis de ambiente para sua sessão:
    ```sh
    export $(grep -v '^#' .env | xargs)
    ```
5.  Execute a aplicação:
    ```sh
    go run ./cmd/api/main.go
    ```

## ⚙️ API / Contrato de Uso

O BFF expõe um único endpoint.

### `POST /weather`

Recebe um CEP em um corpo JSON, valida, e busca as informações de clima.

**Corpo da Requisição:**
```json
{
  "cep": "01001000"
}
```

Exemplos de uso com curl:

- Requisição com Sucesso:

```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"cep": "01001000"}' \
     http://localhost:8081/weather
```

Resposta Esperada (200 OK):

```json
{"city":"São Paulo","temp_C":19.0,"temp_F":66.2,"temp_K":292.0}
```

- CEP com Formato Inválido (Erro de Validação no BFF):

```bash
curl -i -X POST -H "Content-Type: application/json" \
     -d '{"cep": "123"}' \
     http://localhost:8081/weather
```

Resposta Esperada (422 Unprocessable Entity):

```json
{"message":"Request validation failed", ...}
```

CEP Não Encontrado (Erro Repassado do `weather-service`):

```bash
curl -i -X POST -H "Content-Type: application/json" \
     -d '{"cep": "99999999"}' \
     http://localhost:8081/weather
```

Resposta Esperada (404 Not Found):

```json
{"message":"can not find zipcode","code":"not_found"}
```

## 🛠️ Executando os Testes

Para rodar a suíte de testes específica do BFF, que valida a lógica do handler e do serviço de forma isolada, execute o seguinte comando dentro do diretório `weather-bff`:

```sh
go test ./...
```

## 📂 Estrutura do Projeto

Este serviço adota uma estrutura baseada nos princípios da **Arquitetura Hexagonal (Ports & Adapters)** e **Clean Architecture**. O objetivo é isolar o *core* da lógica de negócio de detalhes de infraestrutura, resultando em um sistema altamente desacoplado, testável e fácil de manter.

* `cmd/`: Pontos de entrada da aplicação (`main.go`).
* `internal/`: O coração da aplicação, contendo toda a lógica e abstrações.
    * `di/`: Container de Injeção de Dependência (`uber-go/fx`).
    * `port/`: Interfaces ("Portas") que definem os contratos entre as camadas.
    * `service/`: A lógica de negócio do BFF.
    * `handler/`: Adapters de Entrada (handlers HTTP).
    * `adapter/`: Adapters de Saída (neste caso, o cliente do `weather-service`).
    * `model/`: Entidades e DTOs de entrada/saída.
* `pkg/`: Pacotes de utilidades compartilhadas e genéricas (logger, validator, web helpers).
* `config/`: Lógica de configuração com `Viper`.

## 📝 Variáveis de Ambiente

A aplicação é configurada através de variáveis de ambiente, que são lidas a partir de um arquivo `.env` na raiz do projeto.

### Configuração do Servidor e Logger
| Variável | Descrição | Padrão |
| :--- | :--- | :--- |
| `APP_LOGGER_LEVEL`| Nível mínimo de log (`debug`, `info`, `warn`, `error`). | `info` |
| `APP_SERVER_API_HOST` | Endereço IP em que o servidor irá escutar. | `0.0.0.0` |
| `APP_SERVER_API_PORT` | Porta interna do contêiner em que o servidor irá rodar. | `8080` |
| `APP_SERVER_API_RATE_LIMIT` | Requisições por minuto permitidas por IP/Endpoint. | `100` |
| `APP_SERVER_API_READ_TIMEOUT` | Tempo máximo para ler a requisição inteira. | `5s` |

### Configuração de Clientes
| Variável | Descrição | Padrão |
| :--- | :--- | :--- |
| `APP_CLIENTS_WEATHERSERVICE_URL` | URL base do `weather-service` para comunicação interna. | `http://weather-service:8080`|

*Para a lista completa de variáveis, incluindo timeouts e configurações de CORS, consulte o arquivo `.env.example`.*
