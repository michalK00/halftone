package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michalK00/sg-qr/platform/cloud/aws"
	"log"
	"path"
	"strings"
)

const lifetimeSecs int64 = 60

func DeleteObject(key string) error {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return err
	}
	err = client.DeleteObject(context.Background(), key)
	if err != nil {
		log.Printf("Failed DeleteObject, %v \n", err)
		return err
	}

	return nil
}

func ObjectExists(key string) (bool, error) {
	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return false, err
	}
	_, err = client.HeadObject(context.Background(), key)
	if err != nil {
		log.Printf("Failed ObjectExists, %v \n", err)
		return false, err
	}
	return true, nil
}

func GetObjectUrl(key string) (string, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return "", err
	}
	presignerClient := aws.NewPresignClient(client)
	request, err := presignerClient.GetObjectUrl(context.Background(), key, lifetimeSecs)
	if err != nil {
		log.Printf("Failed GetObjectUrl, %v \n", err)
		return "", nil
	}

	return request.URL, nil
}

func PostObjectRequest(key string, conditions []interface{}) (*s3.PresignedPostRequest, error) {

	client, err := aws.GetAWSClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return nil, err
	}
	presignerClient := aws.NewPresignClient(client)
	request, err := presignerClient.PostObjectRequest(context.Background(), key, lifetimeSecs, conditions)
	if err != nil {
		log.Printf("Failed PutObjectUrl, %v \n", err)
		return nil, err
	}
	fmt.Println(request)

	return request, nil
}

func BuildObjectKey(dirs []string, objectName, extension string) string {

	fullPath := path.Join(append(dirs, objectName)...)

	if extension != "" {
		fullPath += "." + strings.TrimPrefix(extension, ".")
	}

	return fullPath
}
