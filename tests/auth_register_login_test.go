package tests

import (
	"testing"
	"time"

	protos "github.com/SemenShakhray/protos/gen/go/sso"
	"github.com/SemenShakhray/sso/tests/suite"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	respReg, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId)

	respLogin, err := st.AuthClient.Login(ctx, &protos.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, float64(appID), claims["app_id"].(float64))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistratio(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &protos.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	require.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegisterLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	cases := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with both empty",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
				Email:    tc.email,
				Password: tc.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)

		})

	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	cases := []struct {
		name        string
		email       string
		password    string
		appId       int32
		expectedErr string
	}{
		{
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appId:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    randomFakePassword(),
			appId:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with both empty",
			email:       "",
			password:    "",
			appId:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with non-matching password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appId:       appID,
			expectedErr: "invalid email or password",
		},
		{
			name:        "Login without appID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appId:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &protos.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &protos.LoginRequest{
				Email:    tc.email,
				Password: tc.password,
				AppId:    tc.appId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)

		})
	}
}
func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
