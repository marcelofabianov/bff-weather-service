# Guia de Deploy no Google Cloud Run

Este documento detalha o processo passo a passo para publicar o serviço **Weather API** como um contêiner no Google Cloud Run. A estratégia consiste em usar o Google Cloud Build para construir nossa imagem Docker a partir do código-fonte e enviá-la para o Artifact Registry, de onde o Cloud Run fará o deploy.

## Pré-requisitos

Antes de começar, garanta que você tenha:

1.  **Google Cloud SDK (`gcloud` CLI)**: Instalado e configurado na sua máquina. Se não tiver, siga [este guia](https://cloud.google.com/sdk/docs/install).
2.  **Projeto no Google Cloud**: Um projeto ativo com o Faturamento habilitado.
3.  **Permissões**: As permissões de IAM necessárias no seu projeto para gerenciar o Cloud Build, Artifact Registry e Cloud Run (ex: `Owner`, `Editor`, ou papéis específicos como `Cloud Run Admin` e `Artifact Registry Admin`).

---
## Passo a Passo do Deploy

### 1. Autenticação

Abra seu terminal e autentique-se com sua conta do Google Cloud:
```sh
gcloud auth login
```

### 2. Configuração do Projeto

Defina o projeto do Google Cloud que você deseja usar:


```bash
gcloud config set project [SEU_PROJECT_ID]
```

### 3. Habilitar APIs Necessárias

Se for a primeira vez usando esses serviços no projeto, habilite as APIs:

```bash
gcloud services enable run.googleapis.com \
    artifactregistry.googleapis.com \
    cloudbuild.googleapis.com
```

### 4. Criar um Repositório no Artifact Registry

O Artifact Registry é onde nossas imagens Docker serão armazenadas.

```bash
gcloud artifacts repositories create weather-api-repo \
    --repository-format=docker \
    --location=us-central1 \
    --description="Repositório para a Weather API"
```

### 5. Construir e Enviar a Imagem

Este comando usa o Cloud Build para executar as instruções do nosso Dockerfile e enviar a imagem resultante para o Artifact Registry. Execute-o a partir da raiz do seu projeto.

```bash
gcloud builds submit --tag us-central1-docker.pkg.dev/[SEU_PROJECT_ID]/weather-api-repo/weather-api:latest
```

### 6. Fazer o Deploy no Cloud Run

Com a imagem no Artifact Registry, podemos publicá-la como um serviço no Cloud Run.

```bash
gcloud run deploy weather-api-service \
    --image=us-central1-docker.pkg.dev/[SEU_PROJECT_ID]/weather-api-repo/weather-api:latest \
    --platform=managed \
    --region=us-central1 \
    --set-env-vars="APP_CLIENTS_WEATHERAPI_KEY=[SUA_CHAVE_WEATHER_API]" \
    --allow-unauthenticated
```

- `--set-env-vars`: Define a variável de ambiente para a chave da API. Importante:
Substitua `[SUA_CHAVE_WEATHER_API]` pela sua chave real. Para um ambiente de produção, veja a seção de boas práticas abaixo.

- `--allow-unauthenticated`: Torna o endpoint público e acessível pela internet.

## Verificação Pós-Deploy

Use a URL fornecida pelo Cloud Run para testar seu serviço com `curl`:

```bash
curl https://[URL_DO_SEU_SERVICO]/weather/01001000
```

__

## Boas Práticas: Gerenciando a API Key com Secret Manage

Passar segredos diretamente no comando de deploy não é a prática mais segura. O ideal é usar o **Google Secret Manager**.

1. Crie o Segredo:

```bash
gcloud secrets create weather-api-key --replication-policy="automatic"
```

2. Adicione a Versão do Segredo:

```bash
# Substitua [SUA_CHAVE_WEATHER_API] pela sua chave real
printf "[SUA_CHAVE_WEATHER_API]" | gcloud secrets versions add weather-api-key --data-file=-
```

3. Faça o Deploy Referenciando o Segredo:

Agora, use a flag `--set-secrets` em vez de `--set-env-vars`.

```bash
gcloud run deploy weather-api-service \
    --image=us-central1-docker.pkg.dev/[SEU_PROJECT_ID]/weather-api-repo/weather-api:latest \
    --platform=managed \
    --region=us-central1 \
    --set-secrets="APP_CLIENTS_WEATHERAPI_KEY=weather-api-key:latest" \
    --allow-unauthenticated
```

O Cloud Run irá buscar o valor do segredo `weather-api-key` e injetá-lo na variável de ambiente `APP_CLIENTS_WEATHERAPI_KEY` de forma segura.

---

## Atualizando o Serviço

Para publicar uma nova versão do código, simplesmente execute os passos 5 e 6 novamente. O Cloud Build criará uma nova imagem e o Cloud Run criará uma nova revisão do serviço com a imagem atualizada.

## Limpando os Recursos

Para evitar custos, você pode remover os recursos criados:

```bash
# Deletar o serviço do Cloud Run
gcloud run services delete weather-api-service --region=us-central1

# Deletar a imagem do Artifact Registry
gcloud artifacts docker images delete us-central1-docker.pkg.dev/[SEU_PROJECT_ID]/weather-api-repo/weather-api --delete-tags

# Deletar o repositório (opcional)
gcloud artifacts repositories delete weather-api-repo --location=us-central1

# Deletar o segredo (opcional)
gcloud secrets delete weather-api-key
```
