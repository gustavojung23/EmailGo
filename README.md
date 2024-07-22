# Email Sender Project

## Descrição

Este projeto é um serviço de envio de e-mails desenvolvido em Go. Ele utiliza o framework Chi para roteamento, GORM para o mapeamento objeto-relacional, e Testify para testes unitários. Mocks são usados para facilitar o teste das funcionalidades do serviço. Para autenticação, o projeto usa Keycloak, e o banco de dados utilizado é o PostgreSQL. O Docker é usado para containerização do ambiente de desenvolvimento.

## Tecnologias Utilizadas

- **Linguagem**: Go
- **Frameworks**:
  - [Chi](https://github.com/go-chi/chi): Micro framework para criação de APIs HTTP.
  - [GORM](https://gorm.io/): ORM para Go.
  - [Testify](https://github.com/stretchr/testify): Framework para testes unitários.
- **Banco de Dados**: PostgreSQL
- **Autenticação**: Keycloak
- **Containerização**: Docker
- **Outros**:
  - Mocks para simulação de comportamentos durante os testes.