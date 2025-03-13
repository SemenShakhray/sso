package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	protos "github.com/SemenShakhray/protos/gen/go/sso"
	"github.com/SemenShakhray/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	T          *testing.T
	Cfg        *config.Config
	AuthClient protos.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(), grpcAdress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: protos.NewAuthClient(cc),
	}
}

func grpcAdress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
