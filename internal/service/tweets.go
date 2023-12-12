package service

import (
	pbTweets "github.com/Verce11o/yata-protos/gen/go/tweets"
	"github.com/Verce11o/yata/internal/config"
	trace "github.com/Verce11o/yata/internal/lib/metrics/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func MakeTweetsServiceClient(cfg config.Services, tracer *trace.JaegerTracing) pbTweets.TweetsClient {

	// TODO: add retry
	cc, err := grpc.Dial(cfg.Tweets.Addr, grpc.WithTransportCredentials(
		insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(tracer.Provider),
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			)),
	)

	if err != nil {
		log.Fatalf("error while connect to tweets client: %s", err)
	}

	return pbTweets.NewTweetsClient(cc)
}
