package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// Token有效期
//const TokenExpireDuration = time.Hour * 2   //改为从yaml文件中获取

// 用于加密签名的密钥
var myKey = []byte("Francis")

//MyClaims 内嵌的的jwt.StandardClaims 只包含了官方字段
//想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"` //用userID参与生成token
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个自己的声明
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), //过期时间  365 *24 * hour
			Issuer: "reddit", //签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c) //一般默认用 HS256 加密
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(myKey)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)

	// 解析token字符串，获取jwt.Token 以及 信息结构体mc
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return myKey, nil //提供服务端的密钥进行解密
	})
	if err != nil {
		return nil, err
	}

	//if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
	//	return claims, nil
	//}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
