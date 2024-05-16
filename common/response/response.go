package response

import (
	"github.com/gin-gonic/gin"
	"github.com/silentrc/toolbox/common/errConst"
	"net/http"
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Tip  string      `json:"tip"`
}

func NewResponse() *response {
	return &response{}
}

func (res response) SuccessResponse() *response {
	return &response{
		Code: 200,
		Msg:  "success",
	}
}

func (res response) SuccessDataResponse(data interface{}) *response {
	return &response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func (res response) ErrResponse(code int, msg string) *response {
	return &response{
		Code: code,
		Msg:  msg,
	}
}

func (res response) ErrMsgResponse(code int, msg, tip string) *response {
	return &response{
		Code: code,
		Msg:  msg,
		Tip:  tip,
	}
}

// SuccessJson 正确输出
func (res response) SuccessJsonData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, res.SuccessDataResponse(data))
	c.Abort()
}

func (res response) SuccessJson(c *gin.Context) {
	c.JSON(http.StatusOK, res.SuccessResponse())
	c.Abort()
}

// FailJson 错误输出
func (res response) FailJson(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
	c.Abort()
}

func (res response) AuthorizeJson(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	})
	c.Abort()
}

// ErrorJson 错误输出
func (res response) ErrorJson(c *gin.Context, data int) {
	c.JSON(http.StatusBadRequest, res.ErrMsgInterface(data))
	c.Abort()
}

func (res response) ErrMsgInterface(code int) *response {
	switch code {
	case errConst.ParamsErrorCode:
		return &response{
			Code: errConst.ParamsErrorCode,
			Msg:  errConst.ParamsErrorMsg,
		}
	case errConst.PermissionErrorCode:
		return &response{
			Code: errConst.PermissionErrorCode,
			Msg:  errConst.PermissionErrorMsg,
		}
	case errConst.ServerErrorCode:
		return &response{
			Code: errConst.ServerErrorCode,
			Msg:  errConst.ServerErrorMsg,
		}
	case errConst.EmailErrorCode:
		return &response{
			Code: errConst.EmailErrorCode,
			Msg:  errConst.EmailErrorMsg,
		}
	case errConst.SignErrorCode:
		return &response{
			Code: errConst.SignErrorCode,
			Msg:  errConst.SignErrorMsg,
		}
	case errConst.RsaErrorCode:
		return &response{
			Code: errConst.RsaErrorCode,
			Msg:  errConst.RsaErrorMsg,
		}
	case errConst.AccountExistCode:
		return &response{
			Code: errConst.AccountExistCode,
			Msg:  errConst.AccountExistMsg,
		}
	case errConst.OrderErrorCode:
		return &response{
			Code: errConst.OrderErrorCode,
			Msg:  errConst.OrderErrorMsg,
		}
	default:
		return &response{
			Code: errConst.ServerErrorCode,
			Msg:  errConst.ServerErrorMsg,
		}
	}

}

func (r *response) ResponseJson(c *gin.Context) {
	c.JSON(http.StatusOK, r)
	c.Abort()
}
