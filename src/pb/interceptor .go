package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthResponse struct {
	UserID int32  `json:"user_id"`
	Error  string `json:"error,omitempty"`
}

type contextKey string

const userIDKey contextKey = "user_id"

func AuthInterceptor(authServiceURL string) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logrus.Error("No metadata in gPRC request")
			return nil, status.Error(codes.Unauthenticated, "no metadata provided")
		}

		authHeaders := meta.Get("authorization")
		if len(authHeaders) == 0 {
			logrus.Error("metadata has no authorization header")
			return nil, status.Error(codes.Unauthenticated, "metadata has no authorization header. authorization header is required")
		}

		bearer := authHeaders[0]
		if len(bearer) < 7 || bearer[:7] != "Bearer " {
			logrus.Error("invalid auth header foramat")
			return nil, status.Error(codes.Unauthenticated, "invalid auth header foramat")
		}

		token := bearer[7:]

		body, err := json.Marshal(map[string]string{
			"access_token": token,
		})
		if err != nil {
			logrus.Error("json marshalling error")
			return nil, status.Error(codes.Internal, "json marshalling error")
		}

		resp, err := http.Post(authServiceURL+"/validate", "application/json", bytes.NewBuffer(body))
		if err != nil {
			logrus.Errorf("auth http request ended with error, %s", err.Error())
			return nil, status.Errorf(codes.Unavailable, "auth http request ended with error, %s", err.Error())
		}
		defer resp.Body.Close()

		var authResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			logrus.Errorf("failed to decode auth responce, %s", err.Error())
			return nil, status.Errorf(codes.Unauthenticated, "failed to decode auth responce, %s", err.Error())
		}

		if resp.StatusCode != http.StatusOK || authResp.Error != "" {
			logrus.Errorf("Auth service error: %s", authResp.Error)
			return nil, status.Error(codes.Unauthenticated, authResp.Error)
		}

		ctx = context.WithValue(ctx, userIDKey, authResp.UserID)

		return handler(ctx, req)
	}
}
