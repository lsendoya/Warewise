package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/lsendoya/Warewise/pkg/logger"
	"strings"
	"time"
)

type TokenJSON struct {
	AtHash               string   `json:"at_hash"`
	Sub                  string   `json:"sub"`
	CognitoGroups        []string `json:"cognito:groups"`
	EmailVerified        bool     `json:"email_verified"`
	CognitoPreferredRole string   `json:"cognito:preferred_role"`
	Iss                  string   `json:"iss"`
	CognitoUsername      string   `json:"cognito:username"`
	CognitoRoles         []string `json:"cognito:roles"`
	Aud                  string   `json:"aud"`
	TokenUse             string   `json:"token_use"`
	AuthTime             int      `json:"auth_time"`
	Exp                  int      `json:"exp"`
	Iat                  int      `json:"iat"`
	Email                string   `json:"email"`
}

func ValidateToken(token string) (bool, []string, error) {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		msg := "the token is not valid"
		logger.Errorf(msg)
		return false, []string{}, errors.New(msg)
	}

	userInfo, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		msg := "token cannot be decoded"
		logger.Errorf(msg)
		return false, []string{}, err
	}

	var tkj TokenJSON

	errUnmarshall := json.Unmarshal(userInfo, &tkj)
	if errUnmarshall != nil {
		msg := "error parsing token to tokenJSON structure"
		logger.Errorf(msg)
		return false, []string{}, errUnmarshall
	}

	now := time.Now()
	timeExpired := time.Unix(int64(tkj.Exp), 0)

	if timeExpired.Before(now) {
		msg := "token expired"
		logger.Errorf(msg)
		return false, []string{}, errors.New(msg)
	}

	logger.Info("the token is valid")

	return true, tkj.CognitoGroups, nil
}
