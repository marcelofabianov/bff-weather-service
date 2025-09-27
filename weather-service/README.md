# Weather API

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Um microsserviço em Go que, ao receber um CEP brasileiro de 8 dígitos, busca a cidade correspondente e retorna as condições climáticas atuais para aquela localidade.

## ✨ Features e Stack Tecnológica

Este projeto foi construído com foco em boas práticas de arquitetura e resiliência.

* **Linguagem**: Go 1.25+
* **Arquitetura**: Hexagonal (Ports & Adapters) e princípios de Clean Architecture.
* **Injeção de Dependência**: `uber-go/fx` para gerenciar o ciclo de vida e a construção dos componentes da aplicação.
* **Roteamento HTTP**: `chi` para um roteamento leve e poderoso.
* **Configuração**: `viper` para carregar configurações a partir de arquivos `.env` e variáveis de ambiente.
* **Logging**: `slog` para logging estruturado em JSON.
* **Padrões de Resiliência**:
    * **Circuit Breaker**: `sony/gobreaker` para proteger a aplicação contra falhas em cascata de serviços externos.
    * **Retry com Backoff Exponencial**: `jpillora/backoff` para lidar com falhas transitórias de rede ou de APIs.
* **Testes**: `stretchr/testify` para asserções e mocks.
* **Containerização**: Docker e Docker Compose para um ambiente de desenvolvimento e produção consistente.

## 🚀 Como Executar

O projeto é totalmente containerizado, então tudo que você precisa é Docker e Docker Compose.

### Pré-requisitos
* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Passos para a Execução

1.  **Clone o Repositório**
    ```sh
    git clone <URL_DO_SEU_REPOSITORIO>
    cd weather-api
    ```

2.  **Configure o Ambiente**
    Crie o seu arquivo de configuração `.env` a partir do exemplo fornecido.
    ```sh
    cp .env.example .env
    ```
    **Importante:** Abra o arquivo `.env` e insira sua chave de API gratuita do [WeatherAPI.com](https://www.weatherapi.com/) no campo `APP_CLIENTS_WEATHERAPI_KEY`.

3.  **Inicie a Aplicação**
    Use o Docker Compose para construir a imagem e iniciar o contêiner.
    ```sh
    docker compose up --build
    ```
    O terminal exibirá os logs da aplicação. Ao ver a mensagem `"msg":"Starting server"`, a API estará pronta para receber requisições em `http://localhost:8080`.

4.  **Pare a Aplicação**
    Para parar o serviço, pressione `Ctrl+C` no terminal onde ele está rodando.

## ⚙️ Uso da API

A API possui um único endpoint para consulta.

### `GET /weather/{zipcode}`

Retorna a temperatura atual para a cidade correspondente ao CEP informado.

**Exemplos de uso com `curl`:**

* **Requisição com Sucesso (CEP de Goiânia, GO):**
    ```sh
    curl http://localhost:8080/weather/74305460
    ```
    *Resposta Esperada (200 OK):*
    ```json
    {"city": "Goiânia","temp_C":31.3,"temp_F":88.34,"temp_K":304.3}
    ```

* **CEP Não Encontrado:**
    ```sh
    curl -i http://localhost:8080/weather/99999999
    ```
    *Resposta Esperada (404 Not Found):*
    ```json
    {"message":"can not find zipcode","code":"not_found"}
    ```

* **CEP com Formato Inválido:**
    ```sh
    curl -i http://localhost:8080/weather/12345
    ```
    *Resposta Esperada (400 Bad Request):*
    ```json
    {"message":"invalid zipcode","code":"invalid_input"}
    ```
* **Usando httpie:**

```bash
http http://localhost:8080/weather/74305460
HTTP/1.1 200 OK
Cache-Control: no-store, no-cache
Content-Length: 76
Content-Security-Policy: default-src 'none'
Content-Type: application/json; charset=utf-8
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Resource-Policy: same-origin
Date: Sat, 27 Sep 2025 17:47:00 GMT
Permissions-Policy: camera=(), microphone=(), geolocation=()
Referrer-Policy: no-referrer
Strict-Transport-Security: max-age=31536000; includeSubDomains
Vary: Origin
X-Content-Type-Options: nosniff
X-Dns-Prefetch-Control: off
X-Download-Options: noopen
X-Frame-Options: deny
X-Ratelimit-Limit: 100
X-Ratelimit-Remaining: 99
X-Ratelimit-Reset: 1758995220

{
    "city": "Goiânia",
    "temp_C": 32.3,
    "temp_F": 90.13999999999999,
    "temp_K": 305.3
}
```

## 🛠️ Executando os Testes

Para rodar a suíte de testes de unidade e integração, execute o seguinte comando na raiz do projeto:

```sh
go test ./...
```

## 📂 Estrutura do Projeto

Este projeto adota uma estrutura baseada nos princípios da **Arquitetura Hexagonal (Ports & Adapters)** e **Clean Architecture**. O objetivo principal é isolar o *core* da lógica de negócio de detalhes de infraestrutura (como o servidor web e clientes HTTP), resultando em um sistema altamente desacoplado, testável e fácil de manter.

Abaixo está o detalhamento dos principais diretórios e suas responsabilidades:

* `cmd/` - Pontos de Entrada da Aplicação
    * `api/main.go`: Nosso único ponto de entrada. Sua única responsabilidade é inicializar o container de injeção de dependência e executar a aplicação.

* `internal/` - O Coração da Aplicação
    * `di/`: O container de Injeção de Dependência (`uber-go/fx`). É responsável por "montar o quebra-cabeça", construindo e conectando todos os componentes da aplicação (adapters, serviços, handlers, etc.).
    * `port/`: As **"Portas"** do hexágono. Define as **interfaces** que servem como contratos entre o core da aplicação (serviço) e o mundo exterior (handlers, adapters). É a camada de abstração que permite o desacoplamento.
    * `service/`: A implementação da lógica de negócio. Orquestra as funcionalidades da aplicação, seguindo as regras de negócio, e é completamente agnóstico a detalhes de infraestrutura.
    * `handler/`: **Adapters de Entrada**. "Traduzem" requisições do mundo externo (neste caso, requisições HTTP) em chamadas para a camada de serviço.
    * `adapter/`: **Adapters de Saída**. "Traduzem" chamadas da camada de serviço em interações com sistemas externos (neste caso, as chamadas para as APIs ViaCEP e WeatherAPI).
    * `model/`: As entidades de domínio. Representam as estruturas de dados centrais e o comportamento do nosso negócio (ex: `Weather`).

* `pkg/` - Pacotes Compartilhados
    * Contém pacotes de utilidades genéricas e compartilhadas que não possuem dependências do nosso domínio de negócio. São ferramentas que poderiam ser facilmente reutilizadas em outros projetos (ex: `logger`, `web`, `validator`).

* `config/` - Configuração
    * Responsável por carregar, validar e prover as configurações da aplicação para todos os outros módulos, utilizando a biblioteca `viper`.

## 📝 Variáveis de Ambiente

A aplicação é configurada através de variáveis de ambiente, seguindo os princípios do [Twelve-Factor App](https://12factor.net/config). Utilizamos a biblioteca `Viper` para carregar essas variáveis a partir de um arquivo `.env` na raiz do projeto.

### Configuração do Servidor Web
| Variável | Descrição | Padrão |
| :--- | :--- | :--- |
| `APP_SERVER_API_HOST` | Endereço IP em que o servidor irá escutar. | `0.0.0.0` |
| `APP_SERVER_API_PORT` | Porta em que o servidor HTTP irá rodar. | `8080` |
| `APP_SERVER_API_RATE_LIMIT` | Requisições por minuto permitidas por IP/Endpoint. | `100` |
| `APP_SERVER_API_READ_TIMEOUT` | Tempo máximo para ler a requisição inteira. | `5s` |
| `APP_SERVER_API_WRITE_TIMEOUT`| Tempo máximo para escrever a resposta inteira. | `10s` |
| `APP_SERVER_API_IDLE_TIMEOUT` | Tempo máximo que uma conexão pode ficar ociosa. | `120s` |
| `APP_SERVER_API_MAXBODYSIZE` | Tamanho máximo do corpo de uma requisição (em bytes). | `1048576` |

### Configuração de Resiliência
| Variável | Descrição | Padrão |
| :--- | :--- | :--- |
| `APP_RESILIENCE_RETRY_MAX_ATTEMPTS` | Nº máximo de retentativas para chamadas a APIs externas. | `3` |
| `APP_RESILIENCE_RETRY_INITIAL_BACKOFF` | Tempo de espera inicial para a primeira retentativa. | `100ms` |
| `APP_RESILIENCE_RETRY_MAX_BACKOFF` | Tempo máximo de espera entre as retentativas. | `2s` |
| `APP_RESILIENCE_BREAKER_MAX_FAILURES` | Nº de falhas consecutivas para abrir o Circuit Breaker. | `5` |
| `APP_RESILIENCE_BREAKER_TIMEOUT` | Tempo que o Circuit Breaker fica aberto antes de se recuperar. | `30s` |

### Configuração de Clientes Externos
| Variável | Descrição | Padrão |
| :--- | :--- | :--- |
| `APP_CLIENTS_WEATHERAPI_URL` | URL base da WeatherAPI. | `http://api.weatherapi.com/v1` |
| `APP_CLIENTS_WEATHERAPI_KEY` | **(Obrigatório)** Chave de API para autenticar na WeatherAPI.com. | `""` |

*Para outras configurações, como as de CORS, por favor, consulte o arquivo `.env.example` para a lista completa de opções.*

## Deploy do projeto no Google Cloud Run

Para detalhes do deploy confira: [DEPLOY.md](DEPLOY.md)
