package controller

import (
	"fmt"
	_ "fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golibrary/Dao"
	"golibrary/gettime"
	md "golibrary/middleware"
	"golibrary/model"
	"log"
	"net/http"
	"time"
)

// RegisterInfo 用户注册信息
// 注意:注册信息可以使用Gin内部的校验工具或者beego的校验工具进行校验
type RegisterInfo struct {
	//Phone int64  `json:"phone"`
	Name string `json:"name"`
	//Pwd   string `json:"password"`
	//Email string `json:"email"`
}

// RegisterUser 用户注册接口
func RegisterUser(c *gin.Context) {
	//var registerInfo RegisterInfo
	//bindErr := c.BindJSON(&registerInfo)
	//if bindErr == nil {
	//	// 用户注册
	//	err := model.Register(registerInfo.Name, registerInfo.Pwd, registerInfo.Phone, registerInfo.Email)
	//
	//	if err == nil {
	//		c.JSON(http.StatusOK, gin.H{
	//			"status": 0,
	//			"msg":    "success ",
	//			"data":   nil,
	//		})
	//	} else {
	//		c.JSON(http.StatusOK, gin.H{
	//			"status": -1,
	//			"msg":    "注册失败",
	//			"data":   nil,
	//		})
	//	}
	//} else {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": -1,
	//		"msg":    "用户注册解析数据失败" + bindErr.Error(),
	//		"data":   nil,
	//	})
	//}
}

// RegisterPage 对应的注册页面的方法
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Register.html", gin.H{})
}

//对应的登录页面的方法
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{})
}

// HomePage 对应的主页面的方法
func HomePage(c *gin.Context) {
	//claims := c.MustGet("claims").(*myjwt.CustomClaims)
	//if claims != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": 0,
	//		"msg":    "token有效",
	//		"data":   claims,
	//	})
	//}

	recordAll, err := Dao.SelectAllBooksUser()
	if err != nil {
		panic("tab_record select err")
	}
	// fmt.Printf("books:%#v\n", booksList) //切片数组
	// for i, v := range booksList {
	// 	fmt.Println(i, v)
	// }
	if len(*recordAll) > 0 {
		//fmt.Println("shujufanhui")
		c.HTML(http.StatusOK, "Home.html", gin.H{
			"num":     1,
			"data":    recordAll,
			"nowTime": gettime.GetTime(),
		})
	} else {
		fmt.Println("000000")
		c.HTML(http.StatusOK, "Home.html", gin.H{
			"num":     1,
			"data":    nil,
			"nowTime": gettime.GetTime(),
		})
	}
}

func UserRecord(c *gin.Context) {

	claims := c.MustGet("claims").(*md.CustomClaims)
	user := claims.Name
	//fmt.Println(user)
	//book_id := c.PostForm("id")
	//fmt.Println(book_id)
	//var _ int
	////var _ Book_id
	var user_id string
	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", user).Scan(&user_id)
	fmt.Println("开始进入获取记录方法")
	recordAll, err := Dao.SelectAllrecord(user_id)
	if err != nil {
		panic("tab_record select err")
	}
	if len(*recordAll) > 0 {
		fmt.Println("存在记录")
		//fmt.Println("shujufanhui")
		c.HTML(http.StatusOK, "UserRecord.html", gin.H{
			"num":     1,
			"data":    recordAll,
			"nowTime": gettime.GetTime(),
		})
	} else {
		fmt.Println("000000")
		c.HTML(http.StatusOK, "UserRecord.html", gin.H{
			"num":     1,
			"data":    nil,
			"nowTime": gettime.GetTime(),
		})
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	//types := c.DefaultPostForm("type", "post")
	username := c.PostForm("Username")
	password := c.PostForm("Password")
	fmt.Sprintf("username:%s,password:%s", username, password)
	var s string
	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", username).Scan(&s)
	//判断用户是否存在
	//存在输出状态1
	//不存在创建用户，保存密码与用户名
	Bool := IsExist(s)
	if Bool {
		//注册状态
		//State["state"] = 1
		//State["text"] = "此用户已存在！"
		fmt.Println("此用户已存在！")
	} else {
		//用户不存在即添加用户
		AddStruct(username, password)
		//State["state"] = 1
		//State["text"] = "注册成功！"
		c.Redirect(http.StatusMovedPermanently, "/user/login")
	}

	//把状态码和注册状态返回给客户端
	//c.String(http.StatusOK, "%v", State)
	//if err != nil {
	//	fmt.Println("exec failed, ", s)
	//	fmt.Println("该用户不存在")
	//	return
	//}
	//password := c.PostForm("Password")
	// c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
	//c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
	//添加注册成功提示
	//跳转至登录页面
	//c.Redirect(http.StatusMovedPermanently, "/user/login")
}

// LoginResult 登陆结果
type LoginResult struct {
	Token string `json:"token"`
	// 用户模型
	Name string `json:"name"`
	//model.User
}

// Login1 登陆接口 用户名和密码登陆
// name,password
func Login1(c *gin.Context) {
	username := c.PostForm("Username")
	password := c.PostForm("Password")
	var user model.Users
	user.Username = username
	var s string
	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", username).Scan(&s)
	//判断用户是否存在
	//不存在创建用户，保存密码与用户名
	ExistBool := IsExist(s)

	if ExistBool {
		// 登陆逻辑校验
		isPass := IsRight(username, password)
		if isPass {
			// 生成Token
			tokenString := generateToken(c, user)
			// 设置响应头信息的 token
			c.SetCookie("Authorization", tokenString, 60, "/", "127.0.0.1", false, true)
			c.Redirect(http.StatusMovedPermanently, "/user/home")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "验证失败",
				"data":   nil,
			})
		}
	} else {
		fmt.Println("用户不存在")
	}

}

//token生成器
func generateToken(c *gin.Context, user model.Users) string {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := md.NewJWT()

	// 构造用户claims信息(负荷)
	claims := md.CustomClaims{
		user.Username,
		//user.Email,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
			Issuer:    "my-project",                    // 签名颁发者
		},
	}

	// 根据claims生成token对象
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
			"data":   nil,
		})
	}

	log.Println(token)
	// 获取用户相关数据
	data := LoginResult{
		Name:  user.Username,
		Token: token,
	}
	fmt.Println(data)
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 0,
	//	"msg":    "登陆成功",
	//	"data":   data,
	//})
	c.Set("token", token)
	return token
}

// GetDataByTime 测试一个需要认证的接口
func GetDataByTime(c *gin.Context) {
	claims := c.MustGet("claims").(*md.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
		})
	}
}

func Home(c *gin.Context) {
	claims := c.MustGet("claims").(*md.CustomClaims)
	if claims != nil {
		recordAll, err := Dao.SelectAllrecord("1")
		if err != nil {
			panic("tab_record select err")
		}
		// fmt.Printf("books:%#v\n", booksList) //切片数组
		// for i, v := range booksList {
		// 	fmt.Println(i, v)
		// }
		if len(*recordAll) > 0 {
			c.HTML(http.StatusOK, "Home.html", gin.H{
				"num":     1,
				"data":    recordAll,
				"nowTime": gettime.GetTime(),
			})
		} else {
			c.HTML(http.StatusOK, "Home.html", gin.H{
				"num":     1,
				"data":    nil,
				"nowTime": gettime.GetTime(),
			})
		}
	}

}

// 结束操作
func Borrow(c *gin.Context) {
	borrowTime := gettime.GetTime()
	claims := c.MustGet("claims").(*md.CustomClaims)
	user := claims.Name
	fmt.Println(user)
	idStr := c.PostForm("id")
	fmt.Println(idStr)
	var id int
	//var _ Book_id
	var user_id string
	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", user).Scan(&user_id)

	if err := c.ShouldBind(&id); err == nil { //{ID:1}
		fmt.Println("借阅成功")
		//这里写sql语句，也即更新表
		_ = Dao.DB.QueryRow("INSERT INTO tab_record(user_id,book_id,tradingTime) VALUES(?,?,?)", user_id, idStr, borrowTime)
		_ = Dao.DB.QueryRow("UPDATE tab_books SET state=0 WHERE id=?", idStr)
		//tab_name := "tab_books"
		//err := Dao.DeleteRow(&tab_name, _id.ID)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		c.Redirect(http.StatusMovedPermanently, "/user/home")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println("end__________")

}

func ReturnBook(c *gin.Context) {
	claims := c.MustGet("claims").(*md.CustomClaims)
	user := claims.Name
	fmt.Println(user)
	idStr := c.PostForm("id")
	fmt.Println(idStr)
	var _ int
	//var _ Book_id
	var user_id string
	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", user).Scan(&user_id)
	var _id Book_id

	if err := c.ShouldBind(&_id); err == nil { //{ID:1}
		fmt.Println(_id)
		//这里写sql语句，也即更新表
		_ = Dao.DB.QueryRow("DELETE FROM tab_record where user_id=? and book_id=?", user_id, idStr)
		_ = Dao.DB.QueryRow("UPDATE tab_books SET state=1 WHERE id=?", idStr)
		fmt.Println("还书成功")
		//tab_name := "tab_books"
		//err := Dao.DeleteRow(&tab_name, _id.ID)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		c.Redirect(http.StatusMovedPermanently, "/user/home")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//fmt.Println("end__________")

}

// IsExist 判断是否存在用户
func IsExist(user string) bool {
	//如果长度为0说明尚未有用户注册
	if len(user) == 0 {
		return false
	} else {
		return true
	}
}

// IsRight 判断密码是否正确
func IsRight(user string, passwd string) bool {
	var s string
	_ = Dao.DB.QueryRow("SELECT id FROM tab_users WHERE username = ? AND `password` =?", user, passwd).Scan(&s)
	if len(s) == 0 {
		return false
	} else {
		return true
	}
}

func AddStruct(name string, passwd string) {
	err := Dao.DB.QueryRow("INSERT INTO tab_users(username,`password`) VALUES (?,?)", name, passwd)
	fmt.Println(err)
}

//
//import (
//	"fmt"
//	jwtgo "github.com/dgrijalva/jwt-go"
//	"github.com/gin-gonic/gin"
//	"library/Dao"
//
//	//mydomain "library/controller"
//	myjwt "library/middleware"
//	"library/util"
//	"k8s.io/klog"
//	"net/http"
//	"time"
//)
//
//type UserController struct {
//	// service or some to access DB method
//}
//
//func NewUserController() *UserController {
//	controller := UserController{}
//	return &controller
//}
//
//func (ctl *UserController) Login(c *gin.Context) {
//	klog.Infof("login to get a token")
//	username := c.PostForm("Username")
//	password := c.PostForm("Password")
//	fmt.Sprintf("username:%s,password:%s", username, password)
//	var s string
//	_ = Dao.DB.QueryRow("select id from tab_users where username = ?", username).Scan(&s)
//
//	var loginReq LoginReq
//	if err := c.ShouldBindJSON(&loginReq); err == nil {
//		//实际当中需要检查用户名和密码的正确性，这里为了简单起见，hardcode，只要和用户是tom，密码是123456就允许通过
//		// check whether username exists and passwd is matched
//		if loginReq.UserName == "tom" && loginReq.Passwd == "123456" {
//			user := User{}
//			user.UserName = loginReq.UserName
//			user.UserId = 0
//			generateToken(c, user.UserName, 30)
//		} else {
//			c.JSON(http.StatusOK, gin.H{
//				"status": -1,
//				"msg":    "验证失败, 用户不存在或者密码不正确",
//			})
//		}
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"status": -1,
//			"msg":    "json 解析失败." + err.Error(),
//		})
//	}
//}
//
///*
//  此工程为了简单，直接将生成token放在controller中
//  有效时间长度，单位是分钟
//*/
//func generateToken(c *gin.Context, user string, expiredTimeByMinute int64) {
//	j := &myjwt.JWT{
//		[]byte(util.SignKey),
//	}
//	claims := myjwt.CustomClaims{
//		//user.UserId,
//		user,
//		//roleId,
//		jwtgo.StandardClaims{
//			NotBefore: int64(time.Now().Unix() - 1000),                   // 签名生效时间
//			ExpiresAt: int64(time.Now().Unix() + expiredTimeByMinute*60), // 过期时间 一小时
//			Issuer:    "ginjwtdemo",                                      //签名的发行者
//		},
//	}
//
//	token, err := j.CreateToken(claims)
//
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"status": -1,
//			"msg":    err.Error(),
//		})
//		return
//	}
//
//	data := LoginResp{
//		Token: token,
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"status": 0,
//		"msg":    "登录成功！",
//		"data":   data,
//	})
//	return
//}
//
//func (ctl *UserController) CreateOneUser(c *gin.Context) {
//	klog.Infof("create one user")
//	var req UserCreateReq
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	//	user := mydomain.User{UserId: 1, UserName: userName} 会报cannot use promoted field UserCreateReq.UserName in struct literal of type domain.User
//	// 这种太繁琐
//	//user := mydomain.User{}
//	//user.UserName = req.UserName
//	//user.UserId = 0
//
//	user := User{UserId: 0, UserCreateReq: UserCreateReq{UserName: req.UserName}}
//
//	c.JSON(http.StatusOK, gin.H{
//		"result": user,
//		"msg":    "create user successfully",
//	})
//}
//
//func (ctl *UserController) GetAllUsers(c *gin.Context) {
//	claimsFromContext, _ := c.Get(util.Gin_Context_Key)
//	claims := claimsFromContext.(*myjwt.CustomClaims)
//	currentUser := claims.UserName
//	klog.Infof("get all users, loginUser:%q", currentUser)
//	var users []User
//	for i := 0; i < 3; i++ {
//		userName := fmt.Sprintf("tom%d", i)
//		user := User{UserId: 1}
//		user.UserName = userName
//		users = append(users, user)
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"result": users,
//		"count":  len(users),
//	})
//}
//
//func (ctl *UserController) GetOneUser(c *gin.Context) {
//	userId := c.Param("userId")
//	klog.Infof("get one user by id %q", userId)
//}
//
///*
//  // 匹配的url格式:  /usersfind?username=tom&email=test1@163.com
//*/
//func (ctl *UserController) FindUsers(c *gin.Context) {
//	userName := c.DefaultQuery("username", "张三")
//	email := c.Query("email")
//	// 执行实际搜索，这里只是示例
//	c.String(http.StatusOK, "search user by %q %q", userName, email)
//}
//
//func (ctl *UserController) UpdateOneUser(c *gin.Context) {
//	userId := c.Param("userId")
//	klog.Infof("update user by id %q", userId)
//}
//
//func (ctl *UserController) DeleteOneUser(c *gin.Context) {
//	userId := c.Param("userId")
//	klog.Infof("delete user by id %q", userId)
//
//}
