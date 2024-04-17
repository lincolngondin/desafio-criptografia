# Desafio Criptografia
Esse projeto é a solução para o seguinte [desafio backend](https://github.com/backend-br/desafios/blob/master/cryptography/PROBLEM.md), foi feito com o objetivo principal de estudo. O desafio era implementar a criptografia em um serviço de forma transparente para a API e para as camadas de serviço da aplicação.

## Sobre o Projeto
Foi feito na linguagem Go, na versão 1.22 . Foi utilizado o banco de dados sqlite3 para a persistencia de dados.

O projeto é uma API que houve na porta 8080 e expõe os seguintes endpoints:

+ Para criar uma nova entidade. Os campos user_document e credit_card_token serão armazenados já criptografados no banco de dados usando a criptografia AES.
```
POST /users 
{
  "user_document":"123456789",
  "credit_card_token": "987654321",
  "value": 1234
}
```

+ Para obter uma entidade do banco de dados a partir do id. Os dados criptografados anteriormente serão retornados já descriptografados.
```
GET /users/{id}
```

+ Para deletar uma entidade  a partir de seu id.
```
DELETE /users/{id}
```

+ Para atualizar uma entidade a partir do seu id.
```
PUT /users/{id}
{
  "user_document":"1011121314",
  "credit_card_token": "123456789",
  "value": 5678
}
```

## Como executar
Tenha instalado a linguagem [Go](https://go.dev) e o banco de dados [sqlite3](https://www.sqlite.org/). Siga os seguintes passos:
+ Clone esse repositorio.
```
git clone https://github.com/lincolngondin/desafio-criptografia.git
```
+ Entre no diretorio criado
```
cd desafio-criptografia
```
+ Crie o banco de dados sqlite, e leia o arquivo schema.sql da pasta scripts usando o sqlite3.
```
sqlite3 ./data/database.sqlite < ./scripts/schema.sql
```

+ Execute o programa
```
go run ./cmd/app
```
