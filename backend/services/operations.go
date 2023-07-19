package services

import (
	"context"
	"github.com/arik-kfir/greenstar/backend/model"
	"github.com/rueian/rueidis"
	"github.com/secureworks/errors"
)

type OperationsService struct {
	Service
}

func (s *OperationsService) UpdateOperation(ctx context.Context, id string, op model.OperationChanges) (*model.Operation, error) {
	panic(errors.New("not implemented: %s", "UpdateOperation"))
}

func (s *OperationsService) Operation(ctx context.Context, id string) (*model.Operation, error) {
	op := model.Operation{}
	cmd := s.Redis.B().Get().Key("op:" + id).Build()
	if resp := s.Redis.Do(ctx, cmd); resp.Error() != nil {
		if rueidis.IsRedisNil(resp.Error()) {
			return nil, nil
		} else {
			return nil, errors.New("failed to get operation '%s': %w", id, resp.Error())
		}
	} else if err := resp.DecodeJSON(&op); err != nil {
		return nil, errors.New("failed to decode operation: %w", err)
	} else {
		return &op, nil
	}
}
