package util

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/secureworks/errors"
)

func CreateTopic(ctx context.Context, client *pubsub.Client, topicName string, topicConfig *pubsub.TopicConfig) (*pubsub.Topic, error) {
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		return nil, errors.New("failed checking if topic '%s' exists: %w", topicName, err)
	} else if exists {
		return topic, nil
	} else {
		topic, err := client.CreateTopicWithConfig(ctx, topicName, topicConfig)
		if err != nil {
			return nil, errors.New("failed creating topic '%s': %w", topicName, err)
		} else {
			return topic, nil
		}
	}
}

func CreateSubscription(ctx context.Context, client *pubsub.Client, subscriptionName string, subOpts pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	sub := client.Subscription(subscriptionName)
	if exists, err := sub.Exists(ctx); err != nil {
		return nil, errors.New("failed checking if subscription '%s' exists: %w", subscriptionName, err)
	} else if exists {
		return sub, nil
	}

	sub, err := client.CreateSubscription(ctx, subscriptionName, subOpts)
	if err != nil {
		return nil, errors.New("failed creating subscription '%s': %w", subscriptionName, err)
	} else {
		return sub, nil
	}
}
