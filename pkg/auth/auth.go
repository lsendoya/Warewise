package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenJSON struct {
	Sub           string   `json:"sub"`
	CognitoGroups []string `json:"cognito:groups"`
	Iss           string   `json:"iss"`
	ClientID      string   `json:"client_id"`
	EventID       string   `json:"event_id"`
	TokenUse      string   `json:"token_use"`
	Scope         string   `json:"scope"`
	AuthTime      int64    `json:"auth_time"`
	Exp           int64    `json:"exp"`
	Iat           int64    `json:"iat"`
	Jti           string   `json:"jti"`
	Username      string   `json:"username"`
}

func ParseJWT(tokenString string) (TokenJSON, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return TokenJSON{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		jsonClaims, err := json.Marshal(claims)
		if err != nil {
			return TokenJSON{}, err
		}

		var cognitoToken TokenJSON
		err = json.Unmarshal(jsonClaims, &cognitoToken)
		if err != nil {
			return TokenJSON{}, err
		}

		return cognitoToken, nil
	}

	return TokenJSON{}, fmt.Errorf("invalid token claims")
}
