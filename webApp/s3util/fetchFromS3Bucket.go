package s3util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func generatePreSignedURL() string {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("fish-fish-fish-assets"),
		Key:    aws.String("wasm/main.wasm"),
	}, s3.WithPresignExpires(15*time.Minute)) // Expires in 15 minutes

	if err != nil {
		panic("unable to presign request, " + err.Error())
	}

	return req.URL
}

func HandleGetWasmURL(w http.ResponseWriter, r *http.Request) {
	// Configuration

	// Generate the pre signed URL using utility function
	wasmURL := generatePreSignedURL()

	fmt.Println("Pre-signed URL:", wasmURL)
	// Prepare JSON response
	response := map[string]string{"url": wasmURL}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}
