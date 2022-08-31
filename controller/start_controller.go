package controller

import (
	"github.com/gin-gonic/gin"
	"go_m3u8_down/conf"
	"go_m3u8_down/models"
	"go_m3u8_down/models/response"
	"go_m3u8_down/service"
	"net/http"
)

func IndexHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "admin",
	})
}

// StartDownController  首页
func StartDownController(c *gin.Context) {
	var m3u8Model models.M3u8DownModel
	//// 注册绑定参数
	_ = c.ShouldBindQuery(&m3u8Model)

	if m3u8Model.Name == "" {
		response.FailWithMessage("name is null", c)
		return
	}

	if m3u8Model.Url == "" {
		response.FailWithMessage("url is null", c)
		return
	}
	_, err := service.Runs(conf.DownDir, m3u8Model.Url, m3u8Model.Name)
	if err != nil {
		response.FailWithDetailed(err, err.Error(), c)
		return
	}
	response.Ok(c)
}

func AllM3u8Down(c *gin.Context) {
	down := service.AllM3u8Down()
	//marshal, _ := json.Marshal(down)
	//s := string(marshal)
	response.OkWithData(down, c)
}

func DelM3u8(c *gin.Context) {
	hash, _ := c.GetQuery("hash")
	if hash == "" {
		response.Ok(c)
	}
	err := service.DelM3u8ByHash(hash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
}
func ReDownM3u8(c *gin.Context) {
	hash, _ := c.GetQuery("hash")
	if hash == "" {
		response.Ok(c)
	}
	err := service.ReDownM3u8(hash)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
}
