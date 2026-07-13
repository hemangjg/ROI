package grpcapi_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	grpcapi "github.com/ai-finops/ai-finops/services/identity-service/internal/interfaces/grpc"
	pkgauth "github.com/ai-finops/ai-finops/packages/auth"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

func TestExtAuthzServerCheckValidToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	svc := authapp.NewService(nil, nil, privateKey, &privateKey.PublicKey, 15*time.Minute, 24*time.Hour, nil)
	server := grpcapi.NewExtAuthzServer(svc)

	userID := uuid.New().String()
	orgID := uuid.New().String()
	token, err := pkgauth.SignJWT(privateKey, pkgauth.NewAccessClaims(userID, orgID, []string{"admin"}, time.Hour))
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}

	resp, err := server.Check(context.Background(), &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{
						"authorization": "Bearer " + token,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if resp.GetStatus().GetCode() != int32(codes.OK) {
		t.Fatalf("status code = %d, want %d", resp.GetStatus().GetCode(), codes.OK)
	}

	headers := resp.GetOkResponse().GetHeaders()
	if len(headers) != 3 {
		t.Fatalf("header count = %d, want 3", len(headers))
	}
}

func TestExtAuthzServerCheckMissingToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	svc := authapp.NewService(nil, nil, privateKey, &privateKey.PublicKey, 15*time.Minute, 24*time.Hour, nil)
	server := grpcapi.NewExtAuthzServer(svc)

	resp, err := server.Check(context.Background(), &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Request: &authv3.AttributeContext_Request{
				Http: &authv3.AttributeContext_HttpRequest{
					Headers: map[string]string{},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if resp.GetStatus().GetCode() != int32(codes.Unauthenticated) {
		t.Fatalf("status code = %d, want %d", resp.GetStatus().GetCode(), codes.Unauthenticated)
	}
}