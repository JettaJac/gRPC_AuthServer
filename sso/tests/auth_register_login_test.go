package test

import (
	"testing"
	"time"

	"sso/tests/suite"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/jettajac/wb_L0/tree/main/protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParser, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParser.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, claims["appid"].(float64))
	const deltaSecond = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSecond)

}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

//!!! Дописать тесты + что есть, из в,
