package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinController struct {
	C *gin.Context
}

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *GinController) ResponseSuccess(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: "0000",
		Msg:  "success",
		Data: data,
	})
	return
}

// Response setting gin.JSON
func (g *GinController) ResponseMsgDataError(msg string, data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: "9999",
		Msg:  msg,
		Data: data,
	})
	return
}

// Response setting gin.JSON
func (g *GinController) ResponseError(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: "9999",
		Msg:  "error",
		Data: data,
	})
	return
}

// Response setting gin.JSON
func (g *GinController) ResponseMsgError(msg string) {
	g.C.JSON(http.StatusOK, Response{
		Code: "9999",
		Msg:  msg,
		Data: "",
	})
	return
}
