package clients

import (
	pbNotifications "github.com/Verce11o/yata-protos/gen/go/notifications"
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

func MakeNotificationsServiceClient(cfg config.Services, tracer *trace.JaegerTracing, retriesCount int, timeout time.Duration) pbNotifications.NotificationsClient {

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Unavailable),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.Dial(cfg.Notifications.Addr, grpc.WithTransportCredentials(
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
		log.Fatalf("error while connect to notifications client: %s", err)
	}

	return pbNotifications.NewNotificationsClient(cc)
}
