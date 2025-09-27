# Sistema de API de Clima (BFF + Microsserviço)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8.svg)](https://go.dev/)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-v2.0+-2496ED.svg)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Este repositório contém um sistema de microsserviços em Go projetado para consultar informações de clima a partir de um CEP. A arquitetura é composta por dois serviços principais que trabalham em conjunto: um BFF (Backend-for-Frontend) e um serviço de orquestração de negócio.

## 🏛️ Arquitetura do Sistema

O fluxo de uma requisição segue o padrão API Gateway, onde o `weather-bff` atua como a única porta de entrada, validando e repassando a chamada para o `weather-service`, que executa a lógica de negócio principal.

## 🚀 Módulos de Serviço

O sistema é dividido em dois microsserviços independentes, cada um com seu próprio código-fonte e documentação.

### 1. `weather-bff` (Serviço A)

O **Backend-for-Frontend (BFF)** do sistema. É a única porta de entrada pública, responsável por receber as requisições do cliente, realizar validações de formato e atuar como um cliente para o `weather-service`.

> ➡️ **Para mais detalhes, veja o [README.md do weather-bff](./weather-bff/README.md).**

### 2. `weather-service` (Serviço B)

O **microsserviço de negócio principal**. Ele contém a lógica de orquestração para chamar as APIs externas (ViaCEP e WeatherAPI), aplicar regras de resiliência (Circuit Breaker, Retry com Backoff) e formatar a resposta final.

> ➡️ **Para mais detalhes, veja o [README.md do weather-service](./weather-service/README.md).**

## 🏁 Como Executar o Sistema Completo

O projeto é totalmente orquestrado com Docker Compose, facilitando a inicialização de ambos os serviços.

### Pré-requisitos
* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Passos

1.  **Clone este repositório** para sua máquina local.

2.  **Configure os arquivos de ambiente**:
    * Crie o arquivo `weather-service/.env` a partir do `weather-service/.env.example` e insira sua chave da **WeatherAPI**.
    * Crie o arquivo `weather-bff/.env` a partir do `weather-bff/.env.example`.

3.  **Inicie os serviços**:
    A partir do **diretório raiz** (onde o `docker-compose.yml` está localizado), execute:
    ```sh
    docker compose up --build
    ```

4.  **Teste a aplicação**:
    O sistema estará acessível através do BFF na porta `8081`. Use um cliente HTTP para testar:
    ```sh
    curl -X POST -H "Content-Type: application/json" \
         -d '{"cep": "01001000"}' \
         http://localhost:8081/weather
    ```

## 🛠️ Stack Tecnológica Comum

Ambos os serviços foram construídos sobre uma base tecnológica comum, focada em boas práticas:

* **Linguagem**: Go 1.25+
* **Arquitetura**: Hexagonal (Ports & Adapters)
* **Injeção de Dependência**: `uber-go/fx`
* **Roteamento HTTP**: `chi`
* **Configuração**: `viper`
* **Containerização**: Docker
