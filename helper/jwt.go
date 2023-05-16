package helper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/robbert229/jwt"
)

type JwtPayload struct {
	Name string
	Role string
}

var algorithm = jwt.HmacSha256("secreteKey")

func GenerateToken(jwtPayload *JwtPayload) string {
	claims := jwt.NewClaim()
	claims.Set("auth", jwtPayload)
	claims.Set("exp", time.Now().Add(time.Minute*1).Unix())
	token, err := algorithm.Encode(claims)
	if err != nil {
		PanicIfError(err)
	}
	return token
}

func ValidateToken(token string) (JwtPayload, error) {
	claims, err := algorithm.Decode(token)
	if err != nil {
		errMsg := errors.New(err.Error())
		return JwtPayload{}, errMsg
	}
	auth, err := claims.Get("auth")
	if err != nil {
		errMsg := errors.New(err.Error())
		return JwtPayload{}, errMsg
	}
	bit, err := json.Marshal(auth)
	if err != nil {
		errMsg := errors.New(err.Error())
		return JwtPayload{}, errMsg
	}
	var jwtPayload JwtPayload
	err = json.Unmarshal(bit, &jwtPayload)
	if err != nil {
		errMsg := errors.New(err.Error())
		return JwtPayload{}, errMsg
	}
	return jwtPayload, nil
}
