# Backend Engineer Code Challenge - Levee

Através desse desafio, queremos conhecer suas habilidades de desenvolvimento e principalmente como você resolve problemas.


## Sem tempo para realizar o desafio?

Você pode enviar o link de um pull request com uma contribuição sua para qualquer projeto Open Source ou algum projeto próprio que você acredita que demonstre o seu nível técnico e a qualidade do seu código.
Lembre-se que quanto mais código seu pudermos visualizar, melhor será para te avaliarmos. :)

## Sobre o desafio

Temos um subset de dados com algumas vagas de emprego separadas por categoria (arquivo "jobs.txt"). Veja a estrutura abaixo:

```
partnerId|title|categoryId|ExpiresAt
1123|Vendedor|3|21/1/2018
```

Baseando-se nesses dados, queremos que você crie 2 aplicações:


1. A primeira aplicação deverá ser uma API REST para:
  - criar, ativar e listar as vagas;
  - listar a porcentagem e número de vagas ativas por categoria.

2. A segunda aplicação deverá ser responsável por ler o arquivo "jobs.txt", carregar as informações referente às vagas, e servir a aplicação 1.

Ao carregar esses dados essa aplicação deve ser responsável por criptografá-los e armazená-los.

As informações deverão trafegar entre as aplicações criptografadas, consequentemente com a aplicação 1 tendo a habilidade de descriptografá-las. Fica a seu critério escolher a maneira que as aplicações irão se comunicar.

Certifique-se de:
- Que a aplicação seja `idempotente`, para isso você pode considerar que o atributo partnerID de cada vaga é sempre único durante a importação;
- Importar somente vagas onde a data em que irá expirar for maior do que a data da importação;
- Que todas as vagas recém criadas estejam com o status 'draft'.


### Endpoints

Os endpoints disponíveis na aplicação estão listados abaixo. Os classificados como `Protected` devem exigir autenticação, porém o método utilizado fica a seu critério.


| Name       | Method    | URL                  | Protected |
| ---        | ---       | ---                  | :--:      |
| List       | GET       | /jobs                | ✓         |
| Create     | POST      | /jobs                | ✓         |
| Activate   | POST      | /jobs/{:id}/activate | ✓         |
| Percentage | GET       | /category/{:id}      | ✘         |



## Pré-requisitos

O seu teste deve ter um README com os passos necessários para:
- Rodar o projeto;
- Instalar as dependências;
- Rodar os testes automatizados.

### Tecnologia

Ruby ou Go.

### Bônus hints:

- Gostamos de testes;
- Gostamos de Docker;
- Gostamos de CI;
- Message queues são legais;
- Commits estruturados são bem legais;
- Documentação é vida.

## O que está sendo avaliado

Sua capacidade de compreender um problema, desenhar uma solução e aplicá-la.
Queremos ver como você resolve o problema proposto e suas habilidades com as tecnologias propostas.

## Submissão

Para nos enviar seu código, você pode:

- Fazer um fork desse repositório e nos mandar um pull request.
- Nos dar acesso ao seu repositório no github, bitbucket ou gitlab. Adicione o usuário bonigauglitz.
- Se precisar falar com a gente: engenheria@levee.com.br.
