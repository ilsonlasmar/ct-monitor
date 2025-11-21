# ct-monitor

<img width="864" height="388" alt="Screenshot 2025-11-20 at 22 22 30" src="https://github.com/user-attachments/assets/d10ca6ee-c565-429c-9b45-ae28ec54fb51" />



This is a sample template for ct-monitor - Below is a brief explanation of what we have generated for you:

```bash
.
├── cmd/
│   └── functions/
│       ├── consumer/          # Função Lambda Consumer
│       └── producer/          # Função Lambda Producer
├── internal/
│   └── service/               # Lógica de negócio
│       ├── findDomain.go
│       ├── processor.go
│       ├── secrets.go
│       └── sqs.go
├── pkg/                       # Pacotes reutilizáveis
│   ├── awsdb/
│   ├── certspotter/
│   ├── crtsh/
│   ├── ctlog/
│   ├── googlect/
│   ├── logger/
│   ├── request/
│   └── utils/
├── events/
│   └── event.json            # Eventos de teste
├── template.yaml             # Template SAM
├── samconfig.toml           # Configuração do SAM
├── go.mod
├── go.sum
└── Makefile
```


## Requirements

### Production
- [AWS CLI](https://aws.amazon.com/cli/) configurado com permissões de administrador
- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [Docker](https://www.docker.com/community-edition)
- [Go 1.23+](https://golang.org)

### Local development
- [LocalStack](https://localstack.cloud/)
- [awslocal CLI](https://github.com/localstack/awscli-local)
- [samlocal](https://github.com/localstack/aws-sam-cli-local)

## Instalação

### Instalação do SAM CLI
```bash
# macOS
brew install aws-sam-cli

# Windows
choco install aws-sam-cli

# Linux - siga as instruções em:
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
```

### Instalação das Dependências do Projeto

```bash
# Clone o repositório
git clone <repository-url>
cd ct-monitor

# Instale as dependências Go
go mod tidy
```

## Desenvolvimento Local

### Usando LocalStack

#### 1. Instalação do LocalStack e Ferramentas

```bash
# Instalar LocalStack
pip install localstack

# Instalar awslocal
pip install awscli-local

# Instalar samlocal
pip install aws-sam-cli-local
```

#### 2. Iniciando o LocalStack

```bash
# Iniciar LocalStack com os serviços necessários
localstack start -d
```

#### 3. Configuração Local

```bash
# Configurar perfil AWS local
awslocal configure
# AWS Access Key ID: test
# AWS Secret Access Key: test
# Default region name: us-east-1
# Default output format: json
```

#### 4. Deploy Local usando samlocal

```bash
# Build do projeto
samlocal build

# Deploy local
samlocal deploy --guided
```

#### 5. Testando Localmente

```bash
# Listar stacks
awslocal cloudformation list-stacks

# Verificar recursos criados
awslocal dynamodb list-tables
awslocal sqs list-queues
awslocal apigateway get-rest-apis

# Testar API local
curl "http://OUTPUT_URL:4566/Prod/search?domain=example.com" \
  -H "x-api-key: test-key"
```

## Deploy em Produção

### 1. Build da Aplicação

# Com SAM
```bash
sam build
```

### 2. Deploy Inicial

Para o primeiro deploy, use o comando guiado:

```bash
sam deploy --guided
```

### 3. Deploys Subsequentes

Após a configuração inicial, você pode fazer deploy simplesmente com:

```bash
sam deploy
```

### 4. Monitoramento da Stack

```bash
# Verificar status da stack
aws cloudformation describe-stacks --stack-name ct-monitor

# Ver logs das funções Lambda
sam logs -n ProducerFunction --stack-name ct-monitor --tail
sam logs -n ConsumerFunction --stack-name ct-monitor --tail
```
## Configuração

### Variáveis de Ambiente

O projeto utiliza as seguintes variáveis de ambiente configuradas no [template.yaml](template.yaml):

- `SECRET_NAME`: Nome do secret no Secrets Manager contendo a URL da fila SQS

### Parâmetros Configuráveis

No arquivo [samconfig.toml](samconfig.toml), você pode ajustar:

- `region`: Região AWS para deploy
- `stack_name`: Nome da stack CloudFormation
- `s3_prefix`: Prefixo para objetos S3

## API Reference

### Producer Function

**Endpoint**: `GET /search`

**Parâmetros**:
- `domain` (query parameter): Domínio a ser pesquisado nos CT logs

**Headers**:
- `x-api-key`: API Key requerida (obtida após deploy)

**Exemplo**:
```bash
curl -X GET "https://{api-id}.execute-api.{region}.amazonaws.com/Prod/search?domain=example.com" \
  -H "x-api-key: {your-api-key}"
```

## Testes

### Testes Unitários

```bash
# Executar todos os testes
go test ./...

# Executar testes com verbosidade
go test -v ./...

# Executar testes de um pacote específico
go test -v ./pkg/certspotter/
```

### Teste Local das Funções

```bash
# Testar função Producer localmente

am local invoke ProducerFunction --event events/event.json

# Testar função Consumer localmente
sam local invoke ConsumerFunction --event events/sqs-event.json
```

## Links Úteis

- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [Go AWS SDK v2](https://aws.github.io/aws-sdk-go-v2/)
- [Certificate Transparency](https://certificate.transparency.dev/)
