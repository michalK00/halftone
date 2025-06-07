package aws

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	client           *sqs.Client
	defaultQueueURL  string
	defaultQueueName string
}

type SQSMessage struct {
	QueueURL               string
	MessageBody            string
	DelaySeconds           int32
	MessageAttributes      map[string]types.MessageAttributeValue
	MessageGroupId         string
	MessageDeduplicationId string
}

type LambdaPayload struct {
	QueueURL     string
	EventType    string
	Payload      interface{}
	Metadata     map[string]string
	DelaySeconds int32
}

type SQSReceiveMessageParams struct {
	QueueURL              string
	MaxNumberOfMessages   int32
	WaitTimeSeconds       int32
	VisibilityTimeout     int32
	AttributeNames        []types.QueueAttributeName
	MessageAttributeNames []string
}

func NewSQSClient() *SQSClient {
	return &SQSClient{
		defaultQueueURL:  os.Getenv("AWS_SQS_QUEUE_URL"),
		defaultQueueName: os.Getenv("AWS_SQS_QUEUE_NAME"),
	}
}

func (c *SQSClient) Initialize(ctx context.Context, config *Config) error {
	c.client = sqs.NewFromConfig(config.Config)
	return nil
}

func (c *SQSClient) ListQueues(ctx context.Context) (*sqs.ListQueuesOutput, error) {
	return c.client.ListQueues(ctx, &sqs.ListQueuesInput{})
}

func (c *SQSClient) CreateQueue(ctx context.Context, queueName string, attributes map[string]string) (*sqs.CreateQueueOutput, error) {
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}

	if len(attributes) > 0 {
		input.Attributes = attributes
	}

	return c.client.CreateQueue(ctx, input)
}

func (c *SQSClient) DeleteQueue(ctx context.Context, queueURL string) error {
	_, err := c.client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	})
	return err
}

func (c *SQSClient) GetQueueUrl(ctx context.Context, queueName string) (*sqs.GetQueueUrlOutput, error) {
	return c.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
}

func (c *SQSClient) SendMessage(ctx context.Context, message *SQSMessage) (*sqs.SendMessageOutput, error) {
	queueURL := message.QueueURL
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(message.MessageBody),
	}

	if message.DelaySeconds > 0 {
		input.DelaySeconds = message.DelaySeconds
	}

	if len(message.MessageAttributes) > 0 {
		input.MessageAttributes = message.MessageAttributes
	}

	if message.MessageGroupId != "" {
		input.MessageGroupId = aws.String(message.MessageGroupId)
	}

	if message.MessageDeduplicationId != "" {
		input.MessageDeduplicationId = aws.String(message.MessageDeduplicationId)
	}

	return c.client.SendMessage(ctx, input)
}

func (c *SQSClient) SendMessageBatch(ctx context.Context, queueURL string, messages []types.SendMessageBatchRequestEntry) (*sqs.SendMessageBatchOutput, error) {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	return c.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queueURL),
		Entries:  messages,
	})
}

func (c *SQSClient) ReceiveMessage(ctx context.Context, params *SQSReceiveMessageParams) (*sqs.ReceiveMessageOutput, error) {
	queueURL := params.QueueURL
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
	}

	if params.MaxNumberOfMessages > 0 {
		input.MaxNumberOfMessages = params.MaxNumberOfMessages
	}

	if params.WaitTimeSeconds > 0 {
		input.WaitTimeSeconds = params.WaitTimeSeconds
	}

	if params.VisibilityTimeout > 0 {
		input.VisibilityTimeout = params.VisibilityTimeout
	}

	if len(params.AttributeNames) > 0 {
		input.AttributeNames = params.AttributeNames
	}

	if len(params.MessageAttributeNames) > 0 {
		input.MessageAttributeNames = params.MessageAttributeNames
	}

	return c.client.ReceiveMessage(ctx, input)
}

func (c *SQSClient) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func (c *SQSClient) DeleteMessageBatch(ctx context.Context, queueURL string, entries []types.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchOutput, error) {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	return c.client.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(queueURL),
		Entries:  entries,
	})
}

func (c *SQSClient) ChangeMessageVisibility(ctx context.Context, queueURL, receiptHandle string, visibilityTimeout int32) error {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	_, err := c.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(queueURL),
		ReceiptHandle:     aws.String(receiptHandle),
		VisibilityTimeout: visibilityTimeout,
	})
	return err
}

func (c *SQSClient) GetQueueAttributes(ctx context.Context, queueURL string, attributeNames []types.QueueAttributeName) (*sqs.GetQueueAttributesOutput, error) {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	return c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: attributeNames,
	})
}

func (c *SQSClient) SetQueueAttributes(ctx context.Context, queueURL string, attributes map[string]string) error {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	_, err := c.client.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl:   aws.String(queueURL),
		Attributes: attributes,
	})
	return err
}

func (c *SQSClient) PurgeQueue(ctx context.Context, queueURL string) error {
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	_, err := c.client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueURL),
	})
	return err
}

func (c *SQSClient) SendLambdaPayload(ctx context.Context, lambdaPayload *LambdaPayload) (*sqs.SendMessageOutput, error) {
	queueURL := lambdaPayload.QueueURL
	if queueURL == "" {
		queueURL = c.defaultQueueURL
	}

	messageData := map[string]interface{}{
		"eventType": lambdaPayload.EventType,
		"payload":   lambdaPayload.Payload,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if len(lambdaPayload.Metadata) > 0 {
		messageData["metadata"] = lambdaPayload.Metadata
	}

	messageBody, err := json.Marshal(messageData)
	if err != nil {
		return nil, err
	}

	messageAttributes := map[string]types.MessageAttributeValue{
		"eventType": {
			DataType:    aws.String("String"),
			StringValue: aws.String(lambdaPayload.EventType),
		},
	}

	for key, value := range lambdaPayload.Metadata {
		messageAttributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(queueURL),
		MessageBody:       aws.String(string(messageBody)),
		MessageAttributes: messageAttributes,
	}

	if lambdaPayload.DelaySeconds > 0 {
		input.DelaySeconds = lambdaPayload.DelaySeconds
	}

	return c.client.SendMessage(ctx, input)
}
