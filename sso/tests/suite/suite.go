package suite

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/jettajac/wb_L0/tree/main/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}

}
func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
