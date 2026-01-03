package toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// 爬虫工具
type jwtUtils struct {
	SignKey string
}

func (u *utils) NewJwtUtils(signKey string) *jwtUtils {
	return &jwtUtils{
		SignKey: signKey,
	}
}

// JWTAuth 中间件，检查token
func (j *jwtUtils) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorize")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "请求未携带Authorize，无权限访问",
			})
			c.Abort()
			return
		}

		p := j.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := p.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "请求未携带Authorize，无权限访问",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "解析Authorize错误",
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	ErrTokenExpired     = errors.New("token过期")
	ErrTokenNotValidYet = errors.New("token未激活啊")
	ErrTokenMalformed   = errors.New("错误的token")
	ErrTokenInvalid     = errors.New("无法处理此token")
)

type CustomClaims struct {
	UserID  int64  `json:"user_id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

// 新建一个jwt实例
func (j *jwtUtils) NewJWT() *JWT {
	return &JWT{
		[]byte(j.GetSignKey()),
	}
}

// 获取signKey
func (j *jwtUtils) GetSignKey() string {
	return j.SignKey
}

// 这是SignKey
func (j *jwtUtils) SetSignKey(key string) string {
	j.SignKey = key
	return j.SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	iJwtCustomClaims := CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if err == jwt.ErrTokenMalformed {
			return nil, ErrTokenMalformed
		} else if err == jwt.ErrTokenExpired {
			return nil, ErrTokenExpired
		} else if err == jwt.ErrTokenNotValidYet {
			return nil, ErrTokenNotValidYet
		} else {
			return nil, ErrTokenInvalid
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.RegisteredClaims.ExpiresAt.Time = time.Now().Add(1 * time.Hour)
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}
