package ops

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SendMessage(ctx context.Context, sqsClient *sqs.Client, message string, messageGroupId string, uniqueId string) (bool, error) {

	input := &sqs.SendMessageInput{
		QueueUrl:               aws.String("AWS_SQS_URL"),
		MessageBody:            aws.String(message),
		MessageGroupId:         aws.String(messageGroupId),
		MessageDeduplicationId: aws.String(uniqueId),
	}
	result, err := sqsClient.SendMessage(context.TODO(), input)
	fmt.Print(result.MessageId, result)
	if err != nil {
		log.Fatalf("Failed to send message to SQS: %v", err)
		return false, err
	}
	return true, nil
}
