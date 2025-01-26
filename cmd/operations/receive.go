package ops

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// type SqsMessage struct {
// 	MessageId     string
// 	ReceiptHandle string
// 	Body          string
// 	Attributes    map[string]string
// }

func ReceiveMessage(ctx context.Context, sqsClient *sqs.Client) ([]types.Message, error) {
	result, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(os.Getenv("AWS_SQS_URL")),
		MaxNumberOfMessages: 10,
	})
	if err != nil {
		log.Fatalf("Failed to receive messages from SQS: %v", err)
		return result.Messages, err
	}

	// fmt.Println("Messages received:")
	// for i, message := range result.Messages {
	// 	fmt.Printf("Message %d:\n", i+1)
	// 	fmt.Printf("  Message ID: %s\n", *message.MessageId)
	// 	fmt.Printf("  Body: %s\n", *message.Body)
	// 	fmt.Printf("  Receipt Handle: %s\n", *message.ReceiptHandle)
	// 	fmt.Println("  Attributes:")
	// 	for key, value := range message.Attributes {
	// 		fmt.Printf("    %s: %s\n", key, value)
	// 	}
	// 	fmt.Println("  Message Attributes:")
	// 	for key, value := range message.MessageAttributes {
	// 		fmt.Printf("    %s: %s\n", key, *value.StringValue)
	// 	}
	// 	fmt.Println("-----")
	// }
	return result.Messages, nil
}
