package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type Post struct {
	PostID       string `dynamodbav:"PostID"`
	Title        string `dynamodbav:"Title"`
	Author       string `dynamodbav:"Author"`
	PostDate     string `dynamodbav:"PostDate"` // ISO string
	NumReplies   string `dynamodbav:"NumReplies"`
	Link         string `dynamodbav:"Link"`
	FirstMessage string `dynamodbav:"FirstMessage"`

	// DynamoDB keys
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`
}

type PostRecord struct {
	DynamoClient *dynamodb.Client
	TableName    string
}

func (pr *PostRecord) Create(p *Post) error {
	p.PK = fmt.Sprintf("POST#%s", p.PostID)
	p.SK = "POST"

	av, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	_, err = pr.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(pr.TableName),
		Item:      av,
	})
	return err
}

func (pr *PostRecord) Get(postID string) (*Post, error) {
	pk := fmt.Sprintf("POST#%s", postID)

	out, err := pr.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(pr.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: "POST"},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errors.New("post not found")
	}

	var post Post
	err = attributevalue.UnmarshalMap(out.Item, &post)
	return &post, err
}

func (pr *PostRecord) ListPosts(limit int32) ([]Post, error) {
	out, err := pr.DynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String(pr.TableName),
		FilterExpression: aws.String("SK = :skVal"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":skVal": &types.AttributeValueMemberS{Value: "POST"},
		},
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = attributevalue.UnmarshalListOfMaps(out.Items, &posts)
	return posts, err
}
