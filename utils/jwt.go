package utils

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type EasyToken struct {
	Username string
	Uid      int64
	Expires  int64
}

var (
	verifyKey  string
	ErrAbsent  = "token absent"  // 令牌不存在
	ErrInvalid = "token invalid" // 令牌无效
	ErrExpired = "token expired" // 令牌过期
	ErrOther   = "other error"   // 其他错误
)

func init() {
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}
	verifyKey = appConf.String("jwt::token")
}

func (e EasyToken) GetToken() (string, error) {
	/*claims := &jwt.StandardClaims{
		ExpiresAt: e.Expires, //time.Unix(c.ExpiresAt, 0)
		Issuer:    e.Username,
	}*/
	claims := &jwt.MapClaims{
		"uid": e.Uid,
		"iat": time.Now().Unix(),
		"exp": e.Expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(verifyKey))
	if err != nil {
		log.Println(err)
	}
	return tokenString, err
}

func (e EasyToken) ValidateToken(tokenString string) (ret bool, uid int64, err error) {
	if tokenString == "" {
		return false, 0, errors.New(ErrAbsent)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(verifyKey), nil
	})
	if token == nil {
		return false, 0, errors.New(ErrInvalid)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			if v, ok := claims["uid"]; ok {
				//todo float64转化int64
				uid = int64(v.(float64))
				/*switch vv := v.(type) {
				case float64:
					fmt.Println("vv:::::::", vv)
				}*/
			}
		}
		if !ok {
			fmt.Println("ParseHStoken:claims类型转换失败")
			return false, 0, errors.New("ParseHStoken:claims类型转换失败")
		}
		return true, uid, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, 0, errors.New(ErrInvalid)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, 0, errors.New(ErrExpired)
		} else {
			return false, 0, errors.New(ErrOther)
		}
	} else {
		return false, 0, errors.New(ErrOther)
	}
}
