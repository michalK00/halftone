package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoClient struct {
	client          *cognitoidentityprovider.Client
	userPoolID      string
	appClientID     string
	appClientSecret string
	region          string
}

func NewCognitoClient() *CognitoClient {
	return &CognitoClient{
		userPoolID:      os.Getenv("AWS_USER_POOL_ID"),
		appClientID:     os.Getenv("AWS_APP_CLIENT_ID"),
		appClientSecret: os.Getenv("AWS_APP_CLIENT_SECRET"),
		region:          os.Getenv("AWS_REGION"),
	}
}

func (c *CognitoClient) Initialize(ctx context.Context, config *Config) error {
	c.client = cognitoidentityprovider.NewFromConfig(config.Config)
	return nil
}

func (c *CognitoClient) ComputeSecretHash(username string) string {
	if c.appClientSecret == "" {
		return ""
	}

	mac := hmac.New(sha256.New, []byte(c.appClientSecret))
	mac.Write([]byte(username + c.appClientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *CognitoClient) SignUp(ctx context.Context, username, password string, attributes map[string]string) (*cognitoidentityprovider.SignUpOutput, error) {
	var userAttributes []types.AttributeType

	for name, value := range attributes {
		userAttributes = append(userAttributes, types.AttributeType{
			Name:  aws.String(name),
			Value: aws.String(value),
		})
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(c.appClientID),
		Username:       aws.String(username),
		Password:       aws.String(password),
		UserAttributes: userAttributes,
	}

	secretHash := c.ComputeSecretHash(username)
	if secretHash != "" {
		input.SecretHash = aws.String(secretHash)
	}

	return c.client.SignUp(ctx, input)
}

func (c *CognitoClient) InitiateAuth(ctx context.Context, username string, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	authParams := map[string]string{
		"USERNAME": username,
		"PASSWORD": password,
	}

	secretHash := c.ComputeSecretHash(username)
	if secretHash != "" {
		authParams["SECRET_HASH"] = secretHash
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       aws.String(c.appClientID),
		AuthParameters: authParams,
	}

	return c.client.InitiateAuth(ctx, input)
}

func (c *CognitoClient) RefreshToken(ctx context.Context, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(c.appClientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}

	return c.client.InitiateAuth(ctx, input)
}

func (c *CognitoClient) ConfirmSignUp(ctx context.Context, username, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.appClientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
	}

	secretHash := c.ComputeSecretHash(username)
	if secretHash != "" {
		input.SecretHash = aws.String(secretHash)
	}

	return c.client.ConfirmSignUp(ctx, input)
}

func (c *CognitoClient) ForgotPassword(ctx context.Context, username string) (*cognitoidentityprovider.ForgotPasswordOutput, error) {
	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(c.appClientID),
		Username: aws.String(username),
	}

	secretHash := c.ComputeSecretHash(username)
	if secretHash != "" {
		input.SecretHash = aws.String(secretHash)
	}

	return c.client.ForgotPassword(ctx, input)
}

func (c *CognitoClient) ConfirmForgotPassword(ctx context.Context, username, confirmationCode, newPassword string) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.appClientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
	}

	secretHash := c.ComputeSecretHash(username)
	if secretHash != "" {
		input.SecretHash = aws.String(secretHash)
	}

	return c.client.ConfirmForgotPassword(ctx, input)
}

func (c *CognitoClient) GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	return c.client.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	})
}
