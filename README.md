# Shortly: Encurtador de URLs Full-Stack na AWS com CI/CD

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-7B42BC?style=for-the-badge&logo=terraform&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)

## 📖 Sobre o Projeto

**Shortly** é uma aplicação full-stack de encurtamento de URLs construída com foco em performance, escalabilidade e práticas modernas de DevOps. O projeto demonstra um fluxo de trabalho completo, desde o desenvolvimento local com Docker até a implantação totalmente automatizada na AWS com Terraform e GitHub Actions.

Este repositório serve como um projeto de portfólio completo, mostrando competências em desenvolvimento backend, frontend, infraestrutura como código (IaC) e automação de CI/CD.

## ✨ Funcionalidades

* **API de Alta Performance:** Backend em Go para processamento rápido de requisições.
* **Interface Reativa:** Frontend moderno em React com TypeScript.
* **Arquitetura Cloud-Native:** Implantado na AWS utilizando serviços gerenciados e serverless como ECS Fargate e DynamoDB.
* **Infraestrutura como Código:** Toda a infraestrutura da AWS é gerenciada de forma declarativa e versionável com Terraform.
* **Pipeline de CI/CD Automatizado:** O fluxo de trabalho de GitHub Actions automatiza o build, o teste (a ser implementado) e o deploy a cada `push` na branch `main`.

## 🛠️ Stack de Tecnologias

| Categoria      | Tecnologia                               |
| -------------- | ---------------------------------------- |
| **Backend** | Go (Golang)                              |
| **Frontend** | React, TypeScript, Vite                  |
| **Banco de Dados** | Amazon DynamoDB (NoSQL)                  |
| **Cloud** | AWS (ECS Fargate, ECR, ALB, VPC, S3)     |
| **Infra as Code** | Terraform                                |
| **CI/CD** | GitHub Actions                           |
| **Containerização** | Docker, Docker Compose                   |

## 🏗️ Arquitetura

A aplicação é dividida em dois componentes principais: o frontend e o backend, ambos hospedados na AWS. O processo de deploy é totalmente automatizado.

#### Fluxo da Aplicação
Usuário Final --> Route 53 --> ALB (Load Balancer) --> ECS Fargate (Container Go) --> DynamoDB

*O frontend (não implementado neste escopo) seria servido por um S3 com CloudFront.*

#### Fluxo de CI/CD
Desenvolvedor --> git push (branch main) --> GitHub Actions --> (Build & Push Imagem) --> Amazon ECR --> (Update Service) --> Amazon ECS


## 🚀 Rodando Localmente

Para executar este projeto em sua máquina local, siga os passos abaixo.

**Pré-requisitos:**
* [Go](https://go.dev/doc/install) (v1.25+)
* [Node.js](https://nodejs.org/en) (v18+)
* [Docker](https://docs.docker.com/get-docker/) e Docker Compose
* [AWS CLI](https://aws.amazon.com/cli/)
* [Terraform](https://developer.hashicorp.com/terraform/downloads)

**Passos:**

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/SEU_USUARIO/SEU_REPOSITORIO.git](https://github.com/SEU_USUARIO/SEU_REPOSITORIO.git)
    cd SEU_REPOSITORIO
    ```

2.  **Inicie o ambiente de backend:**
    O Docker Compose irá subir o container do backend e um banco de dados DynamoDB local.
    ```bash
    docker-compose up --build
    ```

3.  **Crie a tabela no DynamoDB local:**
    Em um novo terminal, execute o comando abaixo para criar a tabela necessária.
    ```bash
    aws dynamodb create-table \
        --table-name shortly-urls \
        --attribute-definitions AttributeName=ShortCode,AttributeType=S \
        --key-schema AttributeName=ShortCode,KeyType=HASH \
        --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --endpoint-url http://localhost:8000
    ```

4.  **Inicie o frontend:**
    Em outro terminal, navegue para a pasta do frontend, instale as dependências e inicie o servidor de desenvolvimento.
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

5.  Acesse `http://localhost:5173` no seu navegador para ver a aplicação frontend. O backend estará disponível em `http://localhost:8080`.

## ☁️ Implantação na AWS (O Fluxo Automatizado)

Este projeto foi desenhado para ser totalmente automatizado.

#### 1. Infraestrutura com Terraform

A pasta `/terraform` contém toda a definição da nossa infraestrutura na AWS. Ao executar `terraform apply`, os seguintes recursos são criados:
* Uma VPC com subnets públicas e privadas.
* Um Application Load Balancer (ALB) para distribuir o tráfego.
* Um repositório ECR para armazenar a imagem Docker do backend.
* Um cluster ECS com um serviço Fargate para rodar o container do backend.
* Uma tabela no DynamoDB.
* Todas as IAM Roles e Security Groups necessários para a comunicação segura entre os serviços.

#### 2. CI/CD com GitHub Actions

**É assim que a aplicação funciona agora!**

Qualquer `git push` para a branch `main` aciona automaticamente o workflow definido em `.github/workflows/deploy.yml`. Este workflow executa os seguintes passos:

1.  **Autenticação:** Faz login na AWS de forma segura usando os segredos configurados no repositório do GitHub.
2.  **Build da Imagem:** Constrói a imagem Docker da aplicação Go.
3.  **Push para o ECR:** Envia a nova imagem Docker para o repositório Amazon ECR.
4.  **Atualização do Serviço:** Força uma nova implantação (force new deployment) no serviço ECS, que irá parar o container antigo e iniciar um novo com a imagem recém-enviada, tudo isso sem downtime para o usuário (Blue/Green deployment nativo do ECS).
