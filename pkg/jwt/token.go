package jwt

import (
	"qqlx/base/apierr"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/reason"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtConf = &Conf{}

type Conf struct {
	Secret string
	Expire time.Duration
	Issuer string
}

func InitConf() error {
	secret, err := conf.GetJwtSecret()
	if err != nil {
		return err
	}
	expirationTime, err := conf.GetJwtExpirationTime()
	if err != nil {
		return err
	}
	jwtConf = &Conf{
		Secret: secret,
		Expire: expirationTime,
		Issuer: conf.GetJwtIssuer(),
	}
	return nil
}

type MyClaims struct {
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
	*jwt.RegisteredClaims
}

// NewClaims creates a new instance of MyCustomClaims with the given userID and zhName.
func NewClaims(userID int, userName string) *MyClaims {
	now := time.Now()
	return &MyClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    conf.GetJwtIssuer(),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConf.Expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

func (c *MyClaims) GenerateToken() (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	token, err = claims.SignedString([]byte(jwtConf.Secret))
	if err != nil {
		return "", apierr.InternalServer().Set(apierr.JwtErrCode, "failed to generate token", err)
	}
	return token, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	var myCustomClaims MyClaims
	token, err := jwt.ParseWithClaims(tokenString, &myCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConf.Secret), nil
	})
	if err != nil {
		return nil, apierr.InternalServer().Set(apierr.JwtErrCode, "failed to parse token", err)
	}
	if claims, ok := token.Claims.(*MyClaims); ok {
		return claims, nil
	}
	return nil, apierr.InternalServer().Set(apierr.JwtErrCode, "failed to parse token", reason.ErrTokenMode)
}

// GetMyClaims 从gin.Context获取MyCustomClaims
func GetMyClaims(c *gin.Context) (*MyClaims, error) {
	cl, ok := c.Get(constant.AuthMidwareKey)
	if !ok {
		return nil, apierr.Unauthorized().Set(apierr.JwtErrCode, "get claims from context failed", reason.ErrTokenInvalid)
	}
	myCustomClaims, ok := cl.(*MyClaims)
	if !ok {
		return nil, apierr.Unauthorized().Set(apierr.JwtErrCode, "get claims from context failed", reason.ErrTokenMode)
	}
	return myCustomClaims, nil
}
