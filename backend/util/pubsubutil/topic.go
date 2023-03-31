package pubsubutil

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

func CreateDevTopic(ctx context.Context, client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		return nil, fmt.Errorf("failed checking if topic '%s' exists: %w", topicName, err)
	} else if exists {
		log.Ctx(ctx).Info().Str("topic", topicName).Msg("Topic exists")
		return topic, nil
	} else {
		topic, err := client.CreateTopicWithConfig(ctx, topicName, &pubsub.TopicConfig{
			Labels: map[string]string{
				"dev": "true",
				"app": "greenstar",
			},
			MessageStoragePolicy: pubsub.MessageStoragePolicy{AllowedPersistenceRegions: []string{"me-west1"}},
			RetentionDuration:    4 * time.Hour,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create topic '%s': %w", topicName, err)
		} else {
			log.Ctx(ctx).Info().Str("topic", topicName).Msg("Topic created")
			return topic, nil
		}
	}
}
