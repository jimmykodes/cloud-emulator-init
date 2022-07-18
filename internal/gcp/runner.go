package gcp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/grpc/codes"

	"github.com/jimmykodes/cloud-emulator-init/internal/atomicerr"
)

func RunE(ctx context.Context, config *Config) error {
	if config == nil {
		log.Println("no gcp resources to create")
		return nil
	}

	if config.Project == "" {
		return fmt.Errorf("missing gcp project id")
	}

	if len(config.Pubsub) > 0 {
		if emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST"); emulatorHost == "" {
			// just need to make sure this isn't trying to run against real pubsub
			return fmt.Errorf("missing PUBSUB_EMULATOR_HOST env var")
		}
		client, err := pubsub.NewClient(ctx, config.Project)
		if err != nil {
			return err
		}

		var (
			atomicErr atomicerr.Error
			wg        sync.WaitGroup
		)

		for _, topic := range config.Pubsub {
			wg.Add(1)
			go func(topic *Topic) {
				defer wg.Done()
				t, err := client.CreateTopic(ctx, topic.Name)
				if err != nil {
					gErr := &apierror.APIError{}
					if !errors.As(err, &gErr) || gErr.GRPCStatus().Code() != codes.AlreadyExists {
						// error isn't an api error or if it is, the code isn't AlreadyExists
						log.Println("error creating topic", err)
						atomicErr.Append(err)
						return
					}
					log.Println("topic already exists, skipping")
					return
				}
				log.Println("created topic:", t.String())
				for _, subscription := range topic.Subscriptions {
					sub, err := client.CreateSubscription(ctx, subscription, pubsub.SubscriptionConfig{Topic: t})
					if err != nil {
						log.Println("error creating subscription", subscription, "for topic", topic.Name)
						atomicErr.Append(err)
						continue
					}
					log.Println("created subscription", sub.String(), "for topic", topic.Name)
				}
			}(topic)
		}
		wg.Wait()
		if err := atomicErr.Err(); err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Project string   `yaml:"project"`
	Pubsub  []*Topic `yaml:"pubsub"`
}

type Topic struct {
	Name          string   `yaml:"name"`
	Subscriptions []string `yaml:"subscriptions"`
}
