package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	//使用自定义中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//指定静态文件和前端模板的位置
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	//登录业务路由
	v1.POST("/login", controller.LoginHandler)

	//社区分类相关路由
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controller.CommunityHandler)           //获取所有社区分类的列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) //获取某个指定ID的社区分类的内容
		v1.POST("/post", controller.CreatePostHandler)              //发布帖子

		v1.GET("/post/:id", controller.GetPostDetailHandler)       //根据ID查看帖子
		v1.GET("/postList", controller.GetPostListHandler)         //查看帖子列表
		v1.GET("/postListWithOrder", controller.GetPostListRouter) //根据Order参数排序类型 获取帖子列表

		v1.POST("/vote", controller.PostVoteController) //投票
	}

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录用户，通过中间件判断请求头中是否有 有效的JWT
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pages not found",
		})
	})

	return r
}
