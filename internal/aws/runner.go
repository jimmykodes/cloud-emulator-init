package aws

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/jimmykodes/cloud-emulator-init/internal/atomicerr"
)

type Config struct {
	EmulatorURL string   `yaml:"emulatorUrl"`
	Region      string   `yaml:"region"`
	SqsQueues   []string `yaml:"sqs"`
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
		Credentials: credentials.NewStaticCredentials("test", "test", "test"),
		Region:      &config.Region,
		Endpoint:    &config.EmulatorURL,
	})
	if err != nil {
		return err
	}

	log.Println("creating aws resources")
	if len(config.SqsQueues) > 0 {
		client := sqs.New(sess)
		var (
			atomicErr atomicerr.Error
			wg        sync.WaitGroup
		)

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
		wg.Wait()
		if err := atomicErr.Err(); err != nil {
			return err
		}
	}
	log.Println("finished creating aws resources")
	return nil
}
