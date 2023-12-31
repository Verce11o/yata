package clients

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func MakeAuthServiceClient(cfg config.Services, tracer *trace.JaegerTracing, retriesCount int, timeout time.Duration) pbSSO.AuthClient {

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Unavailable),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.Dial(cfg.Auth.Addr, grpc.WithTransportCredentials(
		insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(tracer.Provider),
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		log.Fatalf("error while connect to auth client: %s", err)
	}

	return pbSSO.NewAuthClient(cc)
}
