package s3util

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GeneratePresignedURL(bucketName, objectKey, region string, expiry time.Duration) (string, error) {
	ctx := context.TODO()

	// 1. Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	// 2. Create an S3 client
	client := s3.NewFromConfig(cfg)

	// 3. Create a pre-signer
	presignClient := s3.NewPresignClient(client)

	// 4. Generate a pre-signed URL for GetObject
	presignedGetObject, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(o *s3.PresignOptions) {
		o.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign object: %v", err)
	}

	return presignedGetObject.URL, nil
}
