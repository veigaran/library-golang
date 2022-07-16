package routers

import (
	"golibrary/controller"
	myjwt "golibrary/middleware"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	// "library/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 自定义模板函数
	router.SetFuncMap(template.FuncMap{
		"noImg": func(str string) string { // 没有图片
			if str == "0" {
				return "/img/01.jpg"
			}
			return str
		},
		"stated": func(i int) string { // 是否借出
			if i == 0 {
				return "已借出"
			} else {
				return "未借出"
			}
		},
		"method": func(i int) string { // Record中是借还是还
			if i == 0 {
				return "借书"
			} else {
				return "还书"
			}
		},
		"level": func(i int) string { // 静态文件路径
			level := ""
			for n := 0; n < i; n++ {
				level += "../"
			}
			return level
		},
	})
	// 告诉gin框架模板文件引用的静态文件去哪里找
	router.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	router.LoadHTMLGlob("templates/**/*")

	//router.Use(myjwt.JWTAuth())
	router.GET("/", index)

	adminLogin := router.Group("admin/")
	{
		adminLogin.GET("login", controller.AdminLoginPage)
		adminLogin.POST("login", controller.AdminLogin)
	}

	// 路由组
	adminGroup := router.Group("admin/")
	adminGroup.Use(myjwt.JWTAuth("login"))
	{
		// 首页
		adminGroup.GET("index", controller.IndexHandler)
		// books 图书管理
		adminGroup.GET("books", controller.BooksHander)
		// users 用户管理
		adminGroup.GET("users", controller.UsersHander)
		// record 借阅归还记录
		adminGroup.GET("record", controller.RecordsHander)
		// 查询框 模糊查询
		adminGroup.POST("search", controller.SearchPOST)
	}

	// 嵌套路由组	books
	booksGroup := adminGroup.Group("books")

	{
		// add	增
		booksGroup.GET("/add", controller.AddBooks)       //页面
		booksGroup.POST("/add", controller.AddBooks_post) //数据
		// Delete	删
		booksGroup.POST("/delete", controller.DeleteBooks)
		// change	改
		booksGroup.GET("/bookDetails", controller.DetailsBooks) //书本详情页面
		booksGroup.POST("/bookDetails", controller.ChangeBooks) //修改书

		// query 	查书
	}

	// 嵌套路由组	users
	usersGroup := adminGroup.Group("users")

	{
		// add	增
		// booksGroup.GET("/add", controller.AddBooks)       //页面
		usersGroup.POST("/add", controller.AddUsersHandler) //数据
		// // Delete	删
		usersGroup.POST("/delete", controller.DeleteUsersHandler)
		// // change	改
		usersGroup.POST("/change", controller.ChangeUsersHandler) //修改用户信息

		// query 	查用户
	}
	// 嵌套路由组	record
	recordGroup := adminGroup.Group("record")

	{
		// add	增
		recordGroup.POST("/add", controller.AddRecordHandler) //数据
		// Delete	删
		recordGroup.POST("/delete", controller.DeleteRecordHandler)
		// // change	改
		// recordGroup.POST("/change", controller.ChangeRecordHandler) //修改用户信息

		// query 	借阅归还记录
	}

	//注册登录方法,无需拦截
	readerGroup := router.Group("/user/")
	{
		//注册
		readerGroup.POST("register", controller.Register)
		readerGroup.GET("register", controller.RegisterPage) //加入的页面（注意是get）
		//登录
		readerGroup.POST("login", controller.Login1)
		readerGroup.GET("login", controller.LoginPage) //加入的页面（注意是get）
	}

	// 借、阅、记录查询
	sv1 := router.Group("/user")
	//加载自定义的JWTAuth()中间件,在整个sv1的路由组中都生效
	sv1.Use(myjwt.JWTAuth("login"))
	{
		sv1.POST("Borrow", controller.Borrow)
		sv1.POST("Return", controller.ReturnBook)

		sv1.GET("home", controller.HomePage)
		sv1.POST("home", controller.HomePage)

		sv1.GET("record", controller.UserRecord)
		sv1.POST("record", controller.UserRecord)

		sv1.GET("time", controller.GetDataByTime)
		//sv1.GET("user/home", controller.HomePage) //加入的页面（注意是get）

	}

	//router.POST("/user/Borrow", controller.Borrow)
	//router.POST("/user/Return", controller.ReturnBook)
	//
	//
	//router.GET("/user/home", controller.HomePage)
	//router.POST("/user/home", controller.HomePage)
	//
	//router.POST("/user/home", controller.HomePage)
	//router.GET("/time", myjwt.JWTAuthMiddleware(), test)
	//router.POST("/auth", controller.AuthHandler)
	//router.GET("/home", myjwt.JWTAuthMiddleware(), homeHandler)

	return router
}

//默认页面
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{
		"title": "demo主页",
	})
}

func test(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"username": username},
	})
}
