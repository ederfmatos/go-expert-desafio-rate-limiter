# Rate Limiter com Redis - Projeto

Este projeto implementa um **Rate Limiter** em Go utilizando Redis como backend para armazenar as contagens de requisições. O objetivo do rate limiter é controlar o número de requisições feitas a um servidor, com base em dois critérios: **endereço IP** e **token de acesso**. Além disso, há testes automatizados que validam o funcionamento correto do rate limiter.

## Funcionalidade do Rate Limiter

O Rate Limiter é configurável e pode ser ajustado para limitar o número máximo de requisições por segundo tanto para um **IP específico** quanto para um **token de acesso**. Caso o limite seja excedido, o sistema retorna um código de status HTTP **429** (Too Many Requests).

### Requisitos
- **Limitação por IP**: Limita o número de requisições de um IP específico.
- **Limitação por Token**: Limita o número de requisições de um token específico.
- **Expiração**: Após atingir o limite, a expiração é configurável e o IP/token pode fazer novas requisições após o tempo de expiração.

## Como Executar

### 1. Rodando os Testes

Não é necessário rodar o Redis manualmente, pois estamos utilizando o **Testcontainers** para criar e rodar o Redis automaticamente durante os testes.

#### 1.1. Executando os Testes Automatizados

Para rodar os testes, execute o seguinte comando:

```bash
go test -v
```

Esse comando executará todos os testes definidos no arquivo `main_test.go`, incluindo testes de concorrência e validação do rate limiter.

#### 1.2. Testes Específicos

Os testes são executados de forma paralela e validam cenários como:
- Limitação de requisições por **IP**.
- Limitação de requisições por **Token**.
- **Expiração** após o limite ser atingido.
- **Concorrência** de requisições, com verificação de erros quando o limite é excedido.

Os testes são realizados utilizando o **Testcontainers**, que inicia automaticamente o Redis em um contêiner temporário durante a execução dos testes.

## Testando via cURL

Você pode testar o rate limiter utilizando o `curl`. Aqui estão dois exemplos de como testar:

### 2.1. Teste de Limitação com IP

```bash
curl -X GET http://localhost:8080/ -H "X-Forwarded-For: 192.168.1.1"
```

Este comando simula uma requisição com o IP `192.168.1.1`. Se o limite de requisições para esse IP for atingido, a resposta será um código HTTP **429**.

### 2.2. Teste de Limitação com Token

```bash
curl -X GET http://localhost:8080/ -H "API_KEY: abc123"
```

Este comando simula uma requisição com o token `abc123`. Se o limite de requisições para esse token for atingido, a resposta será um código HTTP **429**.

## Como Funciona o Rate Limiter

1. **IP ou Token**: O rate limiter pode ser configurado para limitar requisições com base no **endereço IP** ou no **token de acesso**.
2. **Expiração**: O limite de requisições é configurável e expira após um tempo específico.
3. **Resposta HTTP 429**: Quando o limite de requisições é atingido, o servidor retorna um código HTTP **429 Too Many Requests** e uma mensagem de erro explicando o motivo.