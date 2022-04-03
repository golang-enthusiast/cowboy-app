package dynamodb

import (
	"fmt"
	"strconv"

	"cowboy-app/internal/domain"
	errors "cowboy-app/internal/error"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	awsDynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	awsErrorResourceInUse = "ResourceInUseException"
)

type cowboyRepo struct {
	db        *awsDynamodb.DynamoDB
	tableName string
}

// NewCowboyRepository - creates a new repository.
func NewCowboyRepository(session *awsSession.Session, tableName string) domain.CowboyRepository {
	// Create a new dynamodb client.
	//
	db := awsDynamodb.New(session)

	// Ensure is table exist.
	//
	cowboyTableMustExist(db, tableName)

	// Return result.
	//
	return &cowboyRepo{
		db:        db,
		tableName: tableName,
	}
}

func cowboyTableMustExist(db *awsDynamodb.DynamoDB, tableName string) {
	// Check if the table exist.
	//
	listTablesOutput, err := db.ListTables(&awsDynamodb.ListTablesInput{})

	// Returned internal server error.
	// The application does not know if the table exists or not.
	// Thus, it cannot query the server, so we panic this error.
	//
	if err != nil {
		panic(err)
	}
	for _, table := range listTablesOutput.TableNames {
		// the table already exists and there is no reason to continue.
		if *table == tableName {
			return
		}
	}

	// Create table.
	//
	_, err = db.CreateTable(&awsDynamodb.CreateTableInput{
		AttributeDefinitions: []*awsDynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(domain.JSONFieldName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*awsDynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(domain.JSONFieldName),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &awsDynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	})
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() != awsErrorResourceInUse {
			panic(aerr)
		}
	}

	// Wait for table.
	_ = db.WaitUntilTableExists(&awsDynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
}

func (r *cowboyRepo) FindByName(name string) (*domain.Cowboy, error) {
	result, err := r.db.GetItem(&awsDynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*awsDynamodb.AttributeValue{
			domain.JSONFieldName: {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.NewErrNotFound(
			fmt.Sprintf("Cowboy not found bu name %s", name),
		)
	}
	item := &domain.Cowboy{}
	errM := dynamodbattribute.UnmarshalMap(result.Item, item)
	if errM != nil {
		return nil, errM
	}
	return item, err
}

func (r *cowboyRepo) List(searchCriteria *domain.CowboySearchCriteria) ([]*domain.Cowboy, error) {
	// Build expression.
	filter := expression.Name(domain.JSONFieldHealth).GreaterThan(expression.Value(searchCriteria.HealthGt))
	if len(searchCriteria.ExcludeName) > 0 {
		filter = filter.And(expression.Name(domain.JSONFieldName).NotEqual(expression.Value(searchCriteria.ExcludeName)))
	}

	expr, err := expression.NewBuilder().WithProjection(r.getProjection()).WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}
	params := &awsDynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(r.tableName),
		Limit:                     aws.Int64(int64(searchCriteria.Limit)),
	}

	// Make the DynamoDB Query API call
	result, err := r.db.Scan(params)
	if err != nil {
		return nil, err
	}
	var items []*domain.Cowboy
	for _, i := range result.Items {
		item := &domain.Cowboy{}
		err = dynamodbattribute.UnmarshalMap(i, item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *cowboyRepo) UpdateHealthPoints(name string, health int32) error {
	updateExpression := fmt.Sprintf(
		"set %s = :%s",
		domain.JSONFieldHealth, // Field name
		domain.JSONFieldHealth, // Param name
	)
	input := &awsDynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*awsDynamodb.AttributeValue{
			fmt.Sprintf(":%s", domain.JSONFieldHealth): {
				N: aws.String(strconv.Itoa(int(health))),
			},
		},
		TableName: aws.String(r.tableName),
		Key: map[string]*awsDynamodb.AttributeValue{
			domain.JSONFieldName: {
				S: aws.String(string(name)),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String(updateExpression),
	}
	_, err := r.db.UpdateItem(input)
	return err
}

func (r *cowboyRepo) getProjection() expression.ProjectionBuilder {
	return expression.NamesList(
		expression.Name(domain.JSONFieldName),
		expression.Name(domain.JSONFieldHealth),
		expression.Name(domain.JSONFieldDamage))
}
