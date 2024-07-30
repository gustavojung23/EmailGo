# Email Sender Project

Este projeto é um serviço de envio de e-mails desenvolvido em Go. Ele utiliza o framework Chi para roteamento, GORM para o mapeamento objeto-relacional, e Testify para testes unitários. Mocks são usados para facilitar o teste das funcionalidades do serviço. Para autenticação, o projeto usa Keycloak e OAuth, utilizando JWT para gerenciamento de tokens. O banco de dados utilizado é o PostgreSQL. O Docker é usado para containerização do ambiente de desenvolvimento. O envio de e-mails é realizado com o Gomail, e as variáveis de ambiente são gerenciadas com Godotenv.

## Tecnologias Utilizadas

- **Linguagem**: Go
- **Frameworks**:
  - [Chi](https://github.com/go-chi/chi): Micro framework para criação de APIs HTTP.
  - [GORM](https://gorm.io/): ORM para Go.
  - [Testify](https://github.com/stretchr/testify): Framework para testes unitários.
- **Bibliotecas**:
  - [JWT](https://github.com/dgrijalva/jwt-go): Biblioteca para criação e verificação de tokens JWT.
  - [Gomail](https://github.com/go-gomail/gomail): Biblioteca para envio de e-mails.
  - [OAuth](https://golang.org/x/oauth2): Biblioteca para OAuth 2.0.
  - [Godotenv](https://github.com/joho/godotenv): Biblioteca para carregar variáveis de ambiente de um arquivo `.env`.
  - [GO Jose](https://github.com/go-jose/go-jose): Biblioteca para lidar com JSON Web Encryption (JWE) e JSON Web Signature (JWS).
- **Banco de Dados**: PostgreSQL
- **Autenticação**: Keycloak e OAuth
- **Containerização**: Docker
- **Outros**:
  - Mocks para simulação de comportamentos durante os testes.