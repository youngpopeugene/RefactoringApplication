package api

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const login = "testuser"

func TestCreateBothTokens(t *testing.T) {
	role := "WORKER"

	tokens, err := createBothTokens(login, role)

	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)

	claimsAccess, err := validateJWT(tokens.AccessToken)
	require.NoError(t, err)
	require.Equal(t, login, claimsAccess["login"])
	require.Equal(t, role, claimsAccess["role"])
	require.Equal(t, "access", claimsAccess["type"])

	claimsRefresh, err := validateJWT(tokens.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, login, claimsRefresh["login"])
	require.Equal(t, role, claimsRefresh["role"])
	require.Equal(t, "refresh", claimsRefresh["type"])
}

func TestCreateAccessToken(t *testing.T) {
	role := "DISPATCHER"

	token, err := createAccessToken(login, role)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := validateJWT(token)
	require.NoError(t, err)
	require.Equal(t, login, claims["login"])
	require.Equal(t, role, claims["role"])
	require.Equal(t, "access", claims["type"])
	require.Greater(t, claims["exp"], time.Now().Unix())
}

func TestCreateRefreshToken(t *testing.T) {
	role := "WORKER"

	token, err := createRefreshToken(login, role)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := validateJWT(token)
	require.NoError(t, err)
	require.Equal(t, login, claims["login"])
	require.Equal(t, role, claims["role"])
	require.Equal(t, "refresh", claims["type"])
	require.Greater(t, claims["exp"], time.Now().Unix())
}

func TestValidateJWT(t *testing.T) {
	payload := map[string]any{
		"login": "testuser",
		"role":  "WORKER",
		"type":  "access",
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	fmt.Println(payload["exp"])

	token, err := generateJWT(payload)
	require.NoError(t, err)

	claims, err := validateJWT(token)
	require.NoError(t, err)
	require.Equal(t, payload["login"], claims["login"])
	require.Equal(t, payload["role"], claims["role"])
	require.Equal(t, payload["type"], claims["type"])
	require.Equal(t, payload["exp"], claims["exp"])

	payload["exp"] = time.Now().Add(-time.Minute).Unix()
	token, err = generateJWT(payload)
	require.NoError(t, err)

	_, err = validateJWT(token)
	require.Error(t, err)
	require.Equal(t, "token expire", err.Error())
}
