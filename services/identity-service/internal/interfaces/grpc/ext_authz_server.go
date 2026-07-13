package grpcapi

import (
	"context"
	"strings"

	authapp "github.com/ai-finops/ai-finops/services/identity-service/internal/application/auth"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// ExtAuthzServer implements Envoy's external authorization gRPC API.
type ExtAuthzServer struct {
	authv3.UnimplementedAuthorizationServer
	svc *authapp.Service
}

// NewExtAuthzServer constructs the Envoy ext_authz adapter.
func NewExtAuthzServer(svc *authapp.Service) *ExtAuthzServer {
	return &ExtAuthzServer{svc: svc}
}

// Check validates JWT access tokens for Envoy before routing to protected services.
func (s *ExtAuthzServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	headers := req.GetAttributes().GetRequest().GetHttp().GetHeaders()
	token := authHeaderValue(headers)

	valid, userID, orgID, roles, err := s.svc.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !valid {
		return deniedResponse(codes.Unauthenticated, typev3.StatusCode_Unauthorized), nil
	}

	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					headerValue("x-user-id", userID),
					headerValue("x-org-id", orgID),
					headerValue("x-user-roles", strings.Join(roles, ",")),
				},
			},
		},
	}, nil
}

func authHeaderValue(headers map[string]string) string {
	for key, value := range headers {
		if strings.EqualFold(key, "authorization") {
			return value
		}
	}
	return ""
}

func deniedResponse(code codes.Code, httpStatus typev3.StatusCode) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(code)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: httpStatus},
			},
		},
	}
}

func headerValue(key, value string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{
		Header: &corev3.HeaderValue{
			Key:   key,
			Value: value,
		},
	}
}