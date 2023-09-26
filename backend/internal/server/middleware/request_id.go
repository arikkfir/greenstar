package middleware

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"github.com/google/uuid"
	"net/http"
)

const RequestIDHeaderName = "X-Request-ID"

type ridKeyType struct{}

var ridKey = &ridKeyType{}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(ridKey)
	if v == nil {
		panic(fmt.Errorf("no request ID found in context"))
	} else if rid, ok := v.(string); ok {
		return rid
	} else {
		panic(fmt.Sprintf("unexpected request ID type '%T' encountered: %+v", v, v))
	}
}

func RequestIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeaderName)
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(RequestIDHeaderName, rid)
			w.Header().Set(RequestIDHeaderName, rid)
		}

		ctxWithRID := context.WithValue(r.Context(), ridKey, rid)
		reqWithRID := r.WithContext(ctxWithRID)

		logger := util.Logger(r.Context()).With("rid", rid)
		reqWithRIDLogger := util.RequestWithLogger(reqWithRID, logger)

		next.ServeHTTP(w, reqWithRIDLogger)
	}
}
