package service

import (
	pbSSO "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/Verce11o/yata/internal/config"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func MakeAuthServiceClient(cfg config.Services, tracer *trace.JaegerTracing) pbSSO.AuthClient {

	// TODO: add retry
	cc, err := grpc.Dial(cfg.Auth.Addr, grpc.WithTransportCredentials(
		insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(tracer.Provider),
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			)),
	)

	if err != nil {
		log.Fatalf("error while connect to auth client: %s", err)
	}

	return pbSSO.NewAuthClient(cc)
}
