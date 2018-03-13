# Backend Engineer Code Challenge - Levee

Implementação do desafio [Backend Engineer Code Challenge - Levee](https://github.com/EmpregoLigado/code-challenge)

## TL;DR;

Executar

```
docker-compose build
docker-compose up -d
```

Acessar o serviço REST na porta 8080

# Serviço REST

Endpoints

| Name       | Method    | URL                  | Protected |
| ---        | ---       | ---                  | :--:      |
| List       | GET       | /jobs                | ✓         |
| Create     | POST      | /jobs                | ✓         |
| Activate   | POST      | /jobs/{:id}/activate | ✓         |
| Percentage | GET       | /category/{:id}      | ✘         |

Para acessar os endpoints protegidos é necessário informar um token JWT no header Authorization.
Para gerar este token é necessário utilizar a mesma chave de segredo informada
durante a inicialização do serviço REST.
Caso esteja utilizando as configurações padrão, esta chave será:

```
aeae42cd8f444313a4f300088713e71c
```

O token precisa obrigatóriamente possuir os seguintes dados:

Header

```
{
  "alg": "HS256",
  "typ": "JWT"
}
```

Payload

```
{
  "exp": 33119884799
}
```

Note que campo ```exp``` é o valor em segundos desde 01/01/1970 UTC, ou o Unix Time como é conhecido. 

Utilize o seguinte header para requisições ao serviço rodando com as configurações padrão:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMTE5ODg0Nzk5fQ.nHpr7MgyjAIA_I1de-6baw3WU_CvCEuGO54p9Rruqx4
```

Este token é válido até 3017, o que torna uma péssima idéia rodar em produção com os valores default ;)


Novos tokens podem ser facilmente gerados na ferramenta [jwt.io](http://jwt.io).


Os endpoints retornam formato as respostas em formato json.

Adicionalmente foram definidas as seguintes variáveis no endpoint List para paginação:

- limit: default 100
- page: default 0
- status: active, draft, any default any



# Dependências

Todas as dependências de bibliotecas estão contidas no repositório, gerenciadas pelo go dep. 
A instalação desta ferramenta apenas é necessária caso seja adicionada alguma
nova dependência.

```
go get -u github.com/golang/dep/cmd/dep
```

Como o uso de diretórios de vendoring foi adicionado ao compilador na versão 1.9 do Go, essa é a versão mínima que consegue compilar este código na forma que está organizado.

# Decisões de arquitetura

Para a comunicação entre os serviços, escolhi o utilizar [Twirp](https://blog.twitch.tv/twirp-a-sweet-new-rpc-framework-for-go-5f2febbf35f).

Essa biblioteca é uma criação da equipe do Twich para implementar
RPC sem a necessidade de SSL. Como o Twirp utiliza como base o Protocol Buffer, é relativamente fácil migrar para gRPC caso necessário.

Os dados são transitados entre os serviços utilizando criptografia de chave compartilhada. O algoritmo é do tipo AES/CBC/PKCS e as chaves devem ser de 16, 24 ou 32 bytes para selecionar respectivamente AES-128, AES-192 ou AES-256.
O vetor de inicialização da versão implementada é sempre zero.

O armazenamento no banco de dados utiliza o mesmo mecânismo de criptografia. Este foi aplicado apenas ao campo title como exemplo simplificado de ofuscação dos dados.


# Deploy
No repositório encontra-se um arquivo docker-compose.yml que visa facilitar a geração de imagens Docker e a execução de um ambiente funcional local.

Este compose tem as seguintes funções: 

1. Compilar o código e gerar as imagens.
1. Inicializar um servidor de banco de dados e construir as tabelas necessárias.
2. Configurar e inicializar  o serviço de dados.
3. Configurar e inicializar o serviço REST.

Os detalhes da compilação podem ser encontrados nas receitas deploy/DataDockerfile e deploy/RestDockerfile.

As imagens finais são nomeadas levee/dataserver e levee/restserver.

# Testes
Foi utilizada a biblioteca padrão da Golang para criar os testes.
Assim para executar os mesmos basta o seguinte comando no diretório raiz do código:

```
go test ./...
```

Estes testes são efetuados apenas em memória. Caso a variavél DB esteja configurada o banco de dados fornecido na mesma será executado para testes.

No arquivo deploy/mysql-init/mysql.sql encontram-se os dados necessários para inicializar a tabela do mysql caso deseje efetuar isto manualmente.

Pode-se inicializar apenas o banco de dados via docker-compose caso não exista um mysql rodando na host.  

```
docker-compose run --service-ports db
```

Para executar os testes utilizando este db via docker:

```
DB="mysql://root:codechallenge@/jobdb" go test ./... -cover
```

# Configurações

Para o servidor de dados são necessárias as seguintes configurações:

- Porta TCP, default PORT:8079
- Chave de criptografia, default KEY:"01020304050607080910111213141516"
- Database, default DB:"mysql://root:codechallenge@/jobdb"
- Bootstrap file, default BOOTSTRAP:"jobs.txt"

É possível utilizar dois tipos de backend para database, memory e mysql.

Para utilizar o modo volátil utilize a string de conexão:

```
memory://
```

Para o servidor REST temos as seguintes configurações:

- Porta TCP, default PORT:8080
- Chave de criptografia, default KEY:"01020304050607080910111213141516"
- Segredo JWT, default SECRET:"aeae42cd8f444313a4f300088713e71c"
- Servidor de dados,  default DATA_HOST:"http://data:8079"

A chave de criptografia deve ser a mesma para os dois serviços.


# Observações
Por não possuir o repositório com o mesmo nome dos include paths, não consegui efetuar os testes no travis, mas segue um travis.yml básico.

Para informar uma host na string de conexão do mysql é necessário fornecer o protocolo também. Exemplo:

```
mysql://root:codechallenge@tcp(db:3306)/jobdb
```

Mais informações sobre esta caracteristica do driver em:
[https://github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql#dsn-data-source-name)