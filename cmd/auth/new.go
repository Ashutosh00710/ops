package auth

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getnoops/ops/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/oidc/v2/pkg/oidc"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to NoOps",
		Long:  `Using SSO login to NoOps`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.GetViper())

			Login(config)
		},
	}

	addFlags(cmd)
	return cmd
}

func Login(config *Config) {
	ctx, cancel := context.WithCancel(context.Background())
	tokenChan := make(chan *oidc.Tokens[*oidc.IDTokenClaims], 1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-sigs
		cancel()
	}()

	server, err := NewServer(ctx, config, tokenChan)
	logging.OnError(err).Fatal("failed to create server")

	select {
	case <-ctx.Done():
		os.Exit(0)
	case token := <-tokenChan:
		if err := server.Shutdown(ctx); err != nil {
			logging.OnError(err).Fatal("failed to shutdown server")
		}

		fmt.Printf("AccessToken: %s\n", token.AccessToken)
		fmt.Printf("RefreshToken: %s\n", token.RefreshToken)
		fmt.Printf("TokenType: %s\n", token.TokenType)
	}
}
