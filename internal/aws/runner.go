package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/jimmykodes/cloud-emulator-init/internal/atomicerr"
)

type Config struct {
	EmulatorURL string   `yaml:"emulatorUrl"`
	Region      string   `yaml:"region"`
	SqsQueues   []string `yaml:"sqs"`
	S3Buckets   []string `yaml:"s3"`
}

func RunE(ctx context.Context, config *Config) error {
	if config == nil {
		log.Println("no aws resources to create")
		return nil
	}
	if config.EmulatorURL == "" {
		return fmt.Errorf("missing aws emulator url")
	}
	if config.Region == "" {
		return fmt.Errorf("missing aws region")
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", "test"),
		S3ForcePathStyle: aws.Bool(true),
		Region:           &config.Region,
		Endpoint:         &config.EmulatorURL,
	})
	if err != nil {
		return err
	}

	var (
		atomicErr atomicerr.Error
		wg        sync.WaitGroup
	)

	log.Println("creating aws resources")
	if len(config.SqsQueues) > 0 {
		client := sqs.New(sess)

		for _, queue := range config.SqsQueues {
			wg.Add(1)
			go func(wg *sync.WaitGroup, queue string) {
				defer wg.Done()
				resp, err := client.CreateQueueWithContext(ctx, &sqs.CreateQueueInput{QueueName: &queue})
				if err != nil {
					log.Println("error creating queue", err)
					atomicErr.Append(err)
					return
				}
				log.Println("created queue:", resp.String())
			}(&wg, queue)
		}
	}
	if len(config.S3Buckets) > 0 {
		client := s3.New(sess)
		for _, bucket := range config.S3Buckets {
			wg.Add(1)
			go func(wg *sync.WaitGroup, bucket string) {
				defer wg.Done()
				resp, err := client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{Bucket: &bucket})
				var awsErr awserr.Error
				if errors.As(err, &awsErr) {
					if awsErr.Code() != "BucketAlreadyOwnedByYou" {
						log.Println("error creating bucket", err)
						atomicErr.Append(err)
					}
				} else if err != nil {
					log.Println("error creating bucket", err)
					atomicErr.Append(err)
					return
				}
				log.Println("created bucket:", resp.String())
			}(&wg, bucket)
		}
	}
	wg.Wait()
	if err := atomicErr.Err(); err != nil {
		return err
	}

	log.Println("finished creating aws resources")
	return nil
}
