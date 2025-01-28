package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	MAX_CONTAINERS       = 3
	MAX_TIMEOUT          = time.Second * 40
	IMAGE_NAME           = "code-exec-engine"
	MAX_PROCESSES  int64 = 130
	SigChan        chan os.Signal
	Running        bool = true
)

func LoadAwsConfig() (aws.Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// queueURL := "https://sqs.ap-south-1.amazonaws.com/058264361019/KCQueue.fifo"
	// message := "Hello, SQS!"
	cfg := aws.Config{
		Region: "ap-south-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		),
	}
	return cfg, err
}

func SqsClient() (*sqs.Client, error) {
	config, err := LoadAwsConfig()
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(config)
	return client, err
}

func RedisClient(ctx context.Context) (*redis.Client, error) {
	// fmt.Println(os.Getenv(""), env.REDIS_URL)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: &tls.Config{},
		DB:        0,
	})
	// ctx := context.Background()
	// Ping the Redis server to check the connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
	return rdb, nil
}
