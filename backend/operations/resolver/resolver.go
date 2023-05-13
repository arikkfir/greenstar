package resolver

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rueian/rueidis"
)

type Resolver struct {
	Redis rueidis.Client
}

func (r *Resolver) HandleMessage(ctx context.Context, msg *pubsub.Message) {
	panic("implement me")
}

//func (r *Resolver) RunOpsSubscriber(ctx context.Context) {
//	createOperationSubscriber := natsutil.Subscriber{
//		NATS:        r.NATS,
//		Subject:     "ops.create",
//		DurableName: "ops",
//		MessageHandler: func(ctx context.Context, msg *nats.Msg) error {
//			op := model.Operation{}
//			if err := json.Unmarshal(msg.Data, &op); err != nil {
//				return fmt.Errorf("failed to decode operation: %w", err)
//			} else if b, err := json.Marshal(op); err != nil {
//				return fmt.Errorf("failed to marshal operation to JSON: %w", err)
//			} else {
//				if op.CreatedAt != "" {
//					err := msg.RespondMsg(&nats.Msg{
//						Header: map[string][]string{},
//						Data:   nil,
//					})
//					return fmt.Errorf("")
//				}
//				if err := r.Redis.Do(ctx, r.Redis.B().Set().Key("op:"+op.ID).Value(string(b)).Build()).Error(); err != nil {
//					return fmt.Errorf("failed to store operation in Redis: %w", err)
//				} else {
//					return nil
//				}
//			}
//		},
//	}
//	if err := createOperationSubscriber.Run(ctx); err != nil {
//		// TODO: differentiate between failure to start subscriber, and failure to process
//		log.Warn().Err(err).Msg("Failed to run create-operation subscriber")
//	}
//
//	//updateOperationSubscriber := natsutil.Subscriber{
//	//	NATS:        r.NATS,
//	//	Subject:     "ops.create",
//	//	DurableName: "ops",
//	//	MessageHandler: func(ctx context.Context, msg *nats.Msg) error {
//	//		op := model.Operation{}
//	//		if err := json.Unmarshal(msg.Data, &op); err != nil {
//	//			return fmt.Errorf("failed to decode operation: %w", err)
//	//		} else if b, err := json.Marshal(op); err != nil {
//	//			return fmt.Errorf("failed to marshal operation to JSON: %w", err)
//	//		} else if err := r.Redis.Do(ctx, r.Redis.B().Set().Key("op:"+op.ID).Value(string(b)).Build()).Error(); err != nil {
//	//			return fmt.Errorf("failed to store operation in Redis: %w", err)
//	//		} else {
//	//			return nil
//	//		}
//	//	},
//	//}
//	//if err := updateOperationSubscriber.Run(ctx); err != nil {
//	//	log.Fatal().Err(err).Msg("Failed to run update-operation subscriber")
//	//}
//}
