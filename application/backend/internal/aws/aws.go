package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/michalK00/halftone/platform/cloud/aws"
	"log"
	"path"
	"strings"
)

const lifetimeSecs int64 = 60 * 10

func DeleteObject(key string) error {

	client, err := aws.GetClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return err
	}
	err = client.S3.DeleteObject(context.Background(), key)
	if err != nil {
		log.Printf("Failed DeleteObject, %v \n", err)
		return err
	}

	return nil
}

func ObjectExists(key string) (bool, error) {
	client, err := aws.GetClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return false, err
	}
	_, err = client.S3.HeadObject(context.Background(), key)
	if err != nil {
		log.Printf("Failed ObjectExists, %v \n", err)
		return false, err
	}
	return true, nil
}

func GetObjectUrl(key string) (string, error) {

	client, err := aws.GetClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return "", err
	}
	request, err := client.S3.GetObjectUrl(context.Background(), key, lifetimeSecs)
	if err != nil {
		log.Printf("Failed GetObjectUrl, %v \n", err)
		return "", nil
	}

	return request.URL, nil
}

func PostObjectRequest(key string, conditions []interface{}) (*s3.PresignedPostRequest, error) {

	client, err := aws.GetClient()
	if err != nil {
		log.Printf("Failed GetAWSClient, %v \n", err)
		return nil, err
	}
	request, err := client.S3.PostObjectRequest(context.Background(), key, lifetimeSecs, conditions)
	if err != nil {
		log.Printf("Failed PutObjectUrl, %v \n", err)
		return nil, err
	}

	return request, nil
}

func BuildObjectKey(dirs []string, objectName, extension string) string {

	fullPath := path.Join(append(dirs, objectName)...)

	if extension != "" {
		fullPath += "." + strings.TrimPrefix(extension, ".")
	}

	return fullPath
}
