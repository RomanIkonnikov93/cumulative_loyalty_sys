package authjwt

import (
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

func EncodeJWT(ID, key string) (string, error) {

	// Create the claims
	var claims = model.MyCustomClaims{
		jwt.RegisteredClaims{
			ID: ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}
