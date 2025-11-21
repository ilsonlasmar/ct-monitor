package ctlog

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type store interface {
	updateItem(ctx context.Context, ctlog Ctlog, ctlogId string) (*dynamodb.UpdateItemOutput, error)
	findAllById(ctx context.Context, ctlogId string) ([]Ctlog, error)
}

type dynamoStore struct {
	dynamo *dynamodb.Client
	log    *zap.SugaredLogger
}

func newStore(log *zap.SugaredLogger, dynamo *dynamodb.Client) dynamoStore {
	return dynamoStore{
		log:    log,
		dynamo: dynamo,
	}
}

func (s dynamoStore) updateItem(ctx context.Context, ctlog Ctlog, ctlogId string) (*dynamodb.UpdateItemOutput, error) {
	result, err := s.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: ctlogId},
			"sk": &types.AttributeValueMemberS{Value: ctlog.CommonName},
		},
		TableName: aws.String("ct-logs"),
		UpdateExpression: aws.String(
			"SET #CommonName = :CommonName, " +
				"#IssuerCaID = :IssuerCaID, " +
				"#IssuerName = :IssuerName, " +
				"#NameValue = :NameValue, " +
				"#EntryTimestamp = :EntryTimestamp, " +
				"#SerialNumber = :SerialNumber, " +
				"#NotBefore = :NotBefore, " +
				"#NotAfter = :NotAfter, " +
				"#Source = :Source",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":CommonName":           &types.AttributeValueMemberS{Value: ctlog.CommonName},
			":IssuerCaID":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ctlog.IssuerCaID)},
			":IssuerName":          &types.AttributeValueMemberS{Value: ctlog.IssuerName},
			":NameValue":           &types.AttributeValueMemberS{Value: ctlog.NameValue},
			":EntryTimestamp":      &types.AttributeValueMemberS{Value: ctlog.EntryTimestamp},
			":SerialNumber":        &types.AttributeValueMemberS{Value: ctlog.SerialNumber},
			":NotBefore": &types.AttributeValueMemberS{Value: ctlog.NotBefore},
			":NotAfter":  &types.AttributeValueMemberS{Value: ctlog.NotAfter},
			":Source":    &types.AttributeValueMemberS{Value: ctlog.Source},
		},
		ExpressionAttributeNames: map[string]string{
			"#CommonName": "CommonName",
			"#IssuerCaID":  "IssuerCaID",
			"#IssuerName":  "IssuerName",
			"#NameValue":   "NameValue",
			"#EntryTimestamp": "EntryTimestamp",
			"#SerialNumber":   "SerialNumber",
			"#NotBefore":  "NotBefore",
			"#NotAfter":   "NotAfter",
			"#Source":     "Source",
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s dynamoStore) findAllById(ctx context.Context, ctlogId string) ([]Ctlog, error) {
	ctlogs := []Ctlog{}

	input := &dynamodb.QueryInput{
		TableName: aws.String("ct-logs"),
		KeyConditions: map[string]types.Condition{
			"pk": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: ctlogId},
				},
			},
		},
	}

	result, err := s.dynamo.Query(ctx, input)
	if err != nil {
		return ctlogs, err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &ctlogs)
	if err != nil {
		return ctlogs, ErrMarshalling
	}

	return ctlogs, nil
}
