package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_m3u8_down/controller"
	"net/http"
	"strconv"
)

//Gin 全局 gin 使用
var Gin *gin.Engine

func HttpServer(port int) {
	Gin = gin.Default()

	//导入模板
	Gin.LoadHTMLGlob("templates/**/*")

	// 静态文件服务器
	Gin.Static("/assets", "./resource/assets")
	Gin.Static("/video", "./resource/video")
	//StaticFS搭建文件服务器（可以加载文件夹及文件）
	Gin.StaticFS("/file", http.Dir("./resource/video"))
	// 允许使用跨域请求
	Gin.Use(CORSMiddleware())
	//注册 接口
	Gin.GET("/", controller.IndexHtml)
	Gin.GET("/all", controller.AllM3u8Down)
	Gin.GET("/start_down", controller.StartDownController)
	Gin.GET("/del_m3u8", controller.DelM3u8)
	Gin.GET("/redown_m3u8", controller.ReDownM3u8)

	//启动运行
	fmt.Printf(` 默认前端运行地址:http://127.0.0.1:%d `, port)
	_ = Gin.Run(":" + strconv.Itoa(port))

}

// CORSMiddleware 解决 跨域
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//允许所有请求跨域
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()
	}
}
