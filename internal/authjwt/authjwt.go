package authjwt

import (
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

func EncodeJWT(ID, key string) (string, error) {

	var claims = model.MyCustomClaims{
		ID,
		jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseJWTWithClaims(token string, cfg config.Config) (string, error) {

	tkn, err := jwt.ParseWithClaims(token, &model.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return cfg.JWTSecretKey, nil
	})
	if claims, ok := tkn.Claims.(*model.MyCustomClaims); ok {
		return claims.Foo, nil
	} else {
		return "", err
	}
}
