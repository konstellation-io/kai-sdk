package auth

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type Auth struct {
	logger       logr.Logger
	authEndpoint string
	clientID     string
	clientSecret string
	realm        string
	username     string
	password     string
	jwt          *gocloak.JWT
}

func New(logger logr.Logger) *Auth {
	user := viper.GetString("minio.client_user")
	password := viper.GetString("minio.client_password")
	authEndpoint := viper.GetString("auth.endpoint")
	realm := viper.GetString("auth.realm")
	clientID := viper.GetString("auth.client")
	clientSecret := viper.GetString("auth.client_secret")

	return &Auth{
		logger:       logger,
		authEndpoint: authEndpoint,
		clientID:     clientID,
		clientSecret: clientSecret,
		realm:        realm,
		username:     user,
		password:     password,
		jwt:          nil,
	}
}

func (a *Auth) GetToken() (*gocloak.JWT, error) {
	client := gocloak.NewClient(a.authEndpoint) // "https://auth.kai-dev.konstellation.io"
	ctx := context.Background()

	if a.jwt != nil {
		token, err := client.RefreshToken(ctx, a.jwt.RefreshToken, a.clientID, a.clientSecret, a.realm)

		if err != nil {
			a.logger.V(2).Info("Couldn't refresh token")
		} else {
			a.jwt = token

			return token, nil
		}
	}

	token, err := client.Login(
		ctx,
		a.clientID,     //"kai-kli-oidc",
		a.clientSecret, //"ymCJInhe6GzrraFdwFyJZAflNvohRQ1I",
		a.realm,        // "konstellation",
		a.username,     // "david",
		a.password,     //"password",
	)
	if err != nil {
		a.logger.Info(fmt.Sprintf("Error getting token: %s", err.Error()))
		return nil, err
	}

	a.jwt = token
	return a.jwt, nil
}
