package jwt

import (
	"github.com/Vitokz/smartWorld_Task/config"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

type jwtStruct struct {
	signingKey string
}

type JwtService interface {
	CreateToken(login, name, role string, id int, timeExpired time.Duration) (string, error)
	GetJWTSecret() string
}

type Claims struct {
	jwt.StandardClaims
	Name  string `json:"name"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

func NewJwtService(cfg *config.Config) *jwtStruct {
	return &jwtStruct{
		signingKey: cfg.Jwt.SigningKey,
	}
}

func (j *jwtStruct) CreateToken(login, name, role string, id int, timeExpired time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(id),
			ExpiresAt: time.Now().Add(timeExpired).Unix(),
		},
		Login: login,
		Name:  name,
		Role:  role,
	})

	return token.SignedString([]byte(j.signingKey))
}

func (j *jwtStruct) GetJWTSecret() string {
	return j.signingKey
}
