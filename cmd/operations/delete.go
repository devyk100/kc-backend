package ops

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func DeleteMessage(ctx context.Context, sqsClient *sqs.Client, receiptHandle string) error {
	url := os.Getenv("AWS_SQS_URL")
	_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: &receiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
}
