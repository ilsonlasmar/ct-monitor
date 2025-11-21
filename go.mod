require (
	github.com/aws/aws-lambda-go v1.36.1
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/config v1.31.20
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.23
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.52.6
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.39.13
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.15
	github.com/google/certificate-transparency-go v1.3.2
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.18.24 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.40.2 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

module github.com/ilsonlasmar/ct-monitor

go 1.23.0

toolchain go1.24.10
