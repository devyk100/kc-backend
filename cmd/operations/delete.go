package ops

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func DeleteMessage(ctx context.Context, sqsClient *sqs.Client, message types.Message) error {
	url := os.Getenv("AWS_SQS_URL")
	_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
}
