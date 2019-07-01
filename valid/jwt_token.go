// Package valid
package valid

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"vilgo/vlog"
)

type JwtToken struct {
	EncryptKey []byte
	DecryptKey []byte
}

func NewJwtToken() *JwtToken {
	return &JwtToken{}
}

func (sel *JwtToken) Generate(stash map[string]interface{}) string {
	claim := jwt.MapClaims{
		"exp":     "",
		"nbf":     "",
		"iat":     "",
		"user_id": "888888888",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	str, err := token.SignedString(sel.EncryptKey)
	if err != nil {
		vlog.LogE("err:%s", err.Error())
		return ""
	}
	return str
}

func (sel *JwtToken) Verify(dt string) (stash map[string]interface{}, yes bool) {
	token, err := jwt.ParseWithClaims(dt, jwt.MapClaims{}, func(tk *jwt.Token) (i interface{}, e error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("not SigningMethodHMAC")
		}
		return []byte(sel.DecryptKey), nil
	})
	if err != nil {
		vlog.LogE("err %s", err.Error())
		return nil, false
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, false
	}
	return map[string]interface{}(claim), true
}
