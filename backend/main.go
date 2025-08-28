package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type URLMapping struct {
	ShortCode string `json:"short_code" dynamodbav:"ShortCode"`
	LongURL   string `json:"long_url" dynamodbav:"LongURL"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

var ddbClient *dynamodb.Client
var tableName string
var domainName string

func generateShortCode(url string) string {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode(req.URL)
	mapping := URLMapping{
		ShortCode: shortCode,
		LongURL:   req.URL,
	}

	item, err := attributevalue.MarshalMap(mapping)
	if err != nil {
		http.Error(w, "Failed to marshal item", http.StatusInternalServerError)
		return
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = ddbClient.PutItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to put item to DynamoDB: %v", err)
		http.Error(w, "Failed to store URL", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		ShortURL: "http://" + domainName + "/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ShortCode": &types.AttributeValueMemberS{Value: shortCode},
		},
	}

	result, err := ddbClient.GetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to get item from DynamoDB: %v", err)
		http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	var mapping URLMapping
	err = attributevalue.UnmarshalMap(result.Item, &mapping)
	if err != nil {
		http.Error(w, "Failed to unmarshal item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, mapping.LongURL, http.StatusFound)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		redirectHandler(w, r)
	case http.MethodPost:
		shortenHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	var ok bool
	tableName, ok = os.LookupEnv("DYNAMODB_TABLE_NAME")
	if !ok {
		log.Fatal("DYNAMODB_TABLE_NAME environment variable not set")
	}

	domainName, ok = os.LookupEnv("DOMAIN_NAME")
	if !ok {
		log.Fatal("DOMAIN_NAME environment variable not set")
	}

	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		awsRegion = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	ddbClient = dynamodb.NewFromConfig(cfg)

	mux := http.NewServeMux()
	mux.Handle("/", corsMiddleware(http.HandlerFunc(mainHandler)))
	mux.HandleFunc("/health", healthCheckHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server starting on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}