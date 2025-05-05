# Money Transfer System

Um sistema simples de microserviços desenvolvido em Golang para realizar transferências de dinheiro internacionais entre contas de diferentes países. Este projeto foi criado com fins educacionais para demonstrar conceitos de microserviços, comunicação via gRPC, e uso de banco de dados in-memory com arquivos JSON.

## Visão Geral

O app permite que um cliente envie uma requisição de transferência de dinheiro entre contas de diferentes países (atualmente suportando apenas Brasil - BRL e Estados Unidos - USD). O sistema é composto por três microserviços que se comunicam via gRPC, com um API Gateway como ponto de entrada para requisições HTTP.

### Arquitetura

A arquitetura do sistema segue o padrão de microserviços, com os seguintes componentes:

- **API Gateway**: Recebe requisições HTTP dos clientes, coordena a comunicação com os outros serviços via gRPC e retorna as respostas. Roda na porta `8080`.
- **Transaction Service**: Valida se a conta do remetente tem saldo suficiente para a transferência. Roda na porta `50051`.
- **Conversion Service**: Converte o valor da transferência da moeda do remetente para a moeda do destinatário. Roda na porta `50052`.

O fluxo de uma transferência é o seguinte:

1. O cliente envia uma requisição HTTP POST para o API Gateway com os dados da transferência.
2. O API Gateway chama o Transaction Service via gRPC para validar o saldo.
3. Se o saldo for suficiente, o API Gateway chama o Conversion Service via gRPC para converter o valor da moeda.
4. O API Gateway retorna a resposta ao cliente com o status da operação e o valor convertido (se aplicável).

### Banco de Dados

O sistema utiliza arquivos JSON como banco de dados in-memory para simplicidade:

- **AccountDB.json**: Armazena informações das contas (ID, saldo em centavos, país).
- **CoinDB.json**: Armazena taxas de conversão entre moedas (BRL e USD).
