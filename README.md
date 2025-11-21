# ct-monitor

<img width="864" height="388" alt="Screenshot 2025-11-20 at 22 22 30" src="https://github.com/user-attachments/assets/d10ca6ee-c565-429c-9b45-ae28ec54fb51" />




This is a sample template for ct-monitor - Below is a brief explanation of what we have generated for you:

```bash
.
├── cmd/
│   └── functions/
│       ├── consumer/          # Lambda Consumer Function
│       └── producer/          # Lambda Producer Function
├── internal/
│   └── service/               # Business logic
│       ├── findDomain.go
│       ├── processor.go
│       ├── secrets.go
│       └── sqs.go
├── pkg/                       # Reusable packages
│   ├── awsdb/
│   ├── certspotter/
│   ├── crtsh/
│   ├── ctlog/
│   ├── googlect/
│   ├── logger/
│   ├── request/
│   └── utils/
├── events/
│   └── event.json            # Test events
├── template.yaml             # SAM Template
├── samconfig.toml           # SAM Configuration
├── go.mod
├── go.sum
└── Makefile
```


## Requirements

### Production
- [AWS CLI](https://aws.amazon.com/cli/)
- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [Docker](https://www.docker.com/community-edition)
- [Go 1.23+](https://golang.org)

### Local development
- [LocalStack](https://localstack.cloud/)
- [awslocal CLI](https://github.com/localstack/awscli-local)
- [samlocal](https://github.com/localstack/aws-sam-cli-local)

## Installation

### SAM CLI Installation
```bash
# macOS
brew install aws-sam-cli

# Windows
choco install aws-sam-cli

# Linux - follow instructions at:
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html
```

### Project Dependencies Installation

```bash
# Clone the repository
git clone <repository-url>
cd ct-monitor

# Install Go dependencies
go mod tidy
```

## Local Development

### Using LocalStack

#### 1. LocalStack and Tools Installation

```bash
# Install LocalStack
pip install localstack

# Install awslocal
pip install awscli-local

# Install samlocal
pip install aws-sam-cli-local
```

#### 2. Starting LocalStack

```bash
# Start LocalStack with required services
localstack start -d
```

#### 3. Local Configuration

```bash
# Configure local AWS profile
awslocal configure
# AWS Access Key ID: test
# AWS Secret Access Key: test
# Default region name: us-east-1
# Default output format: json
```

#### 4. Local Deploy using samlocal

```bash
# Build the project
samlocal build

# Local deploy
samlocal deploy --guided
```

#### 5. Local Testing

```bash
# List stacks
awslocal cloudformation list-stacks

# Check created resources
awslocal dynamodb list-tables
awslocal sqs list-queues
awslocal apigateway get-rest-apis

# Test local API
curl "http://OUTPUT_URL:4566/Prod/search?domain=example.com" \
  -H "x-api-key: test-key"
```

## Production Deploy

### 1. Application Build

# With SAM
```bash
sam build
```

### 2. Initial Deploy

For the first deploy, use the guided command:

```bash
sam deploy --guided
```

### 3. Subsequent Deploys

After initial configuration, you can simply deploy with:

```bash
sam deploy
```

### 4. Stack Monitoring

```bash
# Check stack status
aws cloudformation describe-stacks --stack-name ct-monitor

# View Lambda function logs
sam logs -n ProducerFunction --stack-name ct-monitor --tail
sam logs -n ConsumerFunction --stack-name ct-monitor --tail
```
## Configuration

### Environment Variables

The project uses the following environment variables configured in [template.yaml](template.yaml):

- `SECRET_NAME`: Name of the secret in Secrets Manager containing the SQS queue URL

### Configurable Parameters

In the [samconfig.toml](samconfig.toml) file, you can adjust:

- `region`: AWS region for deployment
- `stack_name`: CloudFormation stack name
- `s3_prefix`: Prefix for S3 objects

## API Reference

### Producer Function

**Endpoint**: `GET /search`

**Parameters**:
- `domain` (query parameter): Domain to be searched in CT logs

**Headers**:
- `x-api-key`: Required API Key (obtained after deploy)

**Example**:
```bash
curl -X GET "https://{api-id}.execute-api.{region}.amazonaws.com/Prod/search?domain=example.com" \
  -H "x-api-key: {your-api-key}"
```

## Tests

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with verbosity
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/certspotter/
```

### Local Function Testing

```bash
# Test Producer function locally
sam local invoke ProducerFunction --event events/event.json

# Test Consumer function locally
sam local invoke ConsumerFunction --event events/sqs-event.json
```

## Useful Links

- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [Go AWS SDK v2](https://aws.github.io/aws-sdk-go-v2/)
- [Certificate Transparency](https://certificate.transparency.dev/)
