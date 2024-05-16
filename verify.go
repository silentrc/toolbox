package toolbox

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"regexp"
)

type verifyUtils struct {
}

// 校验类
func (u *utils) NewVerifyUtils() *verifyUtils {
	return &verifyUtils{}
}

// phone verify
func (v *verifyUtils) VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// 邮箱地址校验
func (v *verifyUtils) VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 绑定请求数据
func (v *verifyUtils) BindRequest(c *gin.Context, value interface{}) (err error) {
	if value == nil {
		return errors.New("value is nil")
	}
	err = c.ShouldBind(value)
	if err == nil && value != nil {
		return
	}
	err = c.ShouldBindJSON(value)
	if err == nil && value != nil {
		return
	}
	err = c.ShouldBindQuery(value)
	if err == nil && value != nil {
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyJson, value)
	return
}
