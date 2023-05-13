package pubsubutil

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func CreateDevSubscription(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subscriptionName string) (*pubsub.Subscription, error) {
	sub := client.Subscription(subscriptionName)
	if exists, err := sub.Exists(ctx); err != nil {
		return nil, fmt.Errorf("failed checking if subscription '%s' exists: %w", subscriptionName, err)
	} else if exists {
		log.Ctx(ctx).Info().Str("subscription", subscriptionName).Msg("Subscription exists")
		return sub, nil
	} else {
		sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
			Topic:             topic,
			AckDeadline:       600 * time.Second,
			RetentionDuration: 4 * time.Hour,
			ExpirationPolicy:  24 * time.Hour,
			Labels: map[string]string{
				"dev": "true",
				"app": "greenstar",
			},
			EnableMessageOrdering: true,
			//DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			//	DeadLetterTopic:     topic.ID(),
			//	MaxDeliveryAttempts: 100,
			//},
			RetryPolicy: &pubsub.RetryPolicy{
				MinimumBackoff: 5 * time.Second,
				MaximumBackoff: 10 * time.Minute,
			},
			EnableExactlyOnceDelivery: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription '%s': %w", subscriptionName, err)
		} else {
			log.Ctx(ctx).Info().Str("subscription", subscriptionName).Msg("Subscription created")
			return sub, nil
		}
	}
}
