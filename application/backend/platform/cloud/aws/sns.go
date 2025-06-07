package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNSClient struct {
	client     *sns.Client
	defaultArn string
}

type SNSMessage struct {
	TopicArn          string
	Message           string
	Subject           string
	MessageAttributes map[string]types.MessageAttributeValue
}

type EmailNotification struct {
	Email    string
	Subject  string
	Message  string
	TopicArn string
	Metadata map[string]string
}

func NewSNSClient() *SNSClient {
	return &SNSClient{
		defaultArn: os.Getenv("AWS_SNS_TOPIC_ARN"),
	}
}

func (c *SNSClient) Initialize(ctx context.Context, config *Config) error {
	c.client = sns.NewFromConfig(config.Config)
	return nil
}

func (c *SNSClient) ListTopics(ctx context.Context) (*sns.ListTopicsOutput, error) {
	return c.client.ListTopics(ctx, &sns.ListTopicsInput{})
}

func (c *SNSClient) CreateTopic(ctx context.Context, name string, attributes map[string]string) (*sns.CreateTopicOutput, error) {
	input := &sns.CreateTopicInput{
		Name: aws.String(name),
	}

	if len(attributes) > 0 {
		input.Attributes = attributes
	}

	return c.client.CreateTopic(ctx, input)
}

func (c *SNSClient) DeleteTopic(ctx context.Context, topicArn string) error {
	_, err := c.client.DeleteTopic(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(topicArn),
	})
	return err
}

func (c *SNSClient) Subscribe(ctx context.Context, topicArn, protocol, endpoint string) (*sns.SubscribeOutput, error) {
	return c.client.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String(protocol),
		Endpoint: aws.String(endpoint),
	})
}

func (c *SNSClient) Unsubscribe(ctx context.Context, subscriptionArn string) error {
	_, err := c.client.Unsubscribe(ctx, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscriptionArn),
	})
	return err
}

func (c *SNSClient) Publish(ctx context.Context, message *SNSMessage) (*sns.PublishOutput, error) {
	topicArn := message.TopicArn
	if topicArn == "" {
		topicArn = c.defaultArn
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message.Message),
	}

	if message.Subject != "" {
		input.Subject = aws.String(message.Subject)
	}

	if len(message.MessageAttributes) > 0 {
		input.MessageAttributes = message.MessageAttributes
	}

	return c.client.Publish(ctx, input)
}

func (c *SNSClient) PublishToTarget(ctx context.Context, targetArn, message string, messageAttributes map[string]types.MessageAttributeValue) (*sns.PublishOutput, error) {
	input := &sns.PublishInput{
		TargetArn: aws.String(targetArn),
		Message:   aws.String(message),
	}

	if len(messageAttributes) > 0 {
		input.MessageAttributes = messageAttributes
	}

	return c.client.Publish(ctx, input)
}

func (c *SNSClient) GetTopicAttributes(ctx context.Context, topicArn string) (*sns.GetTopicAttributesOutput, error) {
	return c.client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicArn),
	})
}

func (c *SNSClient) SetTopicAttributes(ctx context.Context, topicArn, attributeName, attributeValue string) error {
	_, err := c.client.SetTopicAttributes(ctx, &sns.SetTopicAttributesInput{
		TopicArn:       aws.String(topicArn),
		AttributeName:  aws.String(attributeName),
		AttributeValue: aws.String(attributeValue),
	})
	return err
}

func (c *SNSClient) SendEmailNotification(ctx context.Context, notification *EmailNotification) (*sns.PublishOutput, error) {
	topicArn := notification.TopicArn
	if topicArn == "" {
		topicArn = c.defaultArn
	}

	messageAttributes := map[string]types.MessageAttributeValue{
		"email": {
			DataType:    aws.String("String"),
			StringValue: aws.String(notification.Email),
		},
	}

	for key, value := range notification.Metadata {
		messageAttributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	input := &sns.PublishInput{
		TopicArn:          aws.String(topicArn),
		Message:           aws.String(notification.Message),
		Subject:           aws.String(notification.Subject),
		MessageAttributes: messageAttributes,
	}

	return c.client.Publish(ctx, input)
}

func (c *SNSClient) ListSubscriptionsByTopic(ctx context.Context, topicArn string) (*sns.ListSubscriptionsByTopicOutput, error) {
	return c.client.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
}
