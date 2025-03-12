package jwt

import (
	"testing"
	"time"

	"github.com/SemenShakhray/sso/internal/domain/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewJWT(t *testing.T) {
	user := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	app := models.App{
		ID:     12,
		Secret: "tets",
	}

	duration := 10 * time.Minute
	expTime := time.Now().Add(duration).Unix()

	tokenString, err := NewToken(user, app, duration)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, float64(user.ID), claims["uid"])
		assert.Equal(t, user.Email, claims["email"])
		assert.Equal(t, float64(app.ID), claims["app_id"])

		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.InDelta(t, expTime, exp, 1.0)
		assert.True(t, exp > float64(time.Now().Unix()))
	} else {
		t.Error("Invalid token claims")
	}

	wrongSecret := "wrongSecret"
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrongSecret), nil
	})

	assert.Error(t, err)
}
