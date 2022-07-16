package controller

import (
	"fmt"
	"golibrary/Dao"
	"golibrary/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	username := c.PostForm("Username")
	password := c.PostForm("Password")
	var user model.Users
	user.Username = username
	var s string
	_ = Dao.DB.QueryRow("SELECT admin_id FROM tab_admin WHERE adminname =? ", username).Scan(&s)
	fmt.Println(s)
	//判断管理员账号是否存在
	ExistBool := IsExist(s)
	if ExistBool {
		// 登陆逻辑校验
		isPass := adminIsright(username, password)
		if isPass {
			// 生成Token
			tokenString := generateToken(c, user)
			// 设置响应头信息的 token
			c.SetCookie("Authorization", tokenString, 60, "/", "127.0.0.1", false, true)
			c.Redirect(http.StatusMovedPermanently, "/admin/users")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "验证失败",
				"data":   nil,
			})
		}
	} else {
		fmt.Println("管理员账号不存在")
		c.Redirect(http.StatusMovedPermanently, "/admin/login")
	}
}

//对应的登录页面的方法
func AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "loginAdmin.html", gin.H{})
}

// users 管理员管理用户的页面
func UsersHander(c *gin.Context) {

	// c.HTML(http.StatusOK, "admin/users.html", nil)

	usersAll, err := Dao.SelectAllUsers()
	if err != nil {
		fmt.Printf("UsersHander err:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err,
		})
	}
	fmt.Printf("users:%v\n", usersAll)
	if len(*usersAll) > 0 {
		c.HTML(http.StatusOK, "admin/users.html", gin.H{
			"num":  1, // http://localhost:8080/admin/users 路径的admin/后面有几级(level)
			"data": *usersAll,
		})
	} else {
		c.HTML(http.StatusOK, "admin/users.html", gin.H{
			"num":  1,
			"data": nil,
		})
	}
}

// 修改用户信息
func ChangeUsersHandler(c *gin.Context) {
	fmt.Printf("c:%#v,%T\n", c, c)
	// var u Dao.Users
	// if err := c.ShouldBind(&u); err == nil {
	idStr := c.PostForm("id")
	flag := c.PostForm("flag")
	value := c.PostForm("value")
	// string到int
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("修改用户信息 err:%v\n", err)
	}
	// 更新book数据
	err = Dao.UpdateRowUser(idInt, &flag, &value)
	if err != nil {
		fmt.Printf("修改用户信息 数据库 err:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"result": 0})

}

// 删除用户信息
func DeleteUsersHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	// string到int
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("删除用户信息 err:%v\n", err)
	}
	// 删除用户数据
	tab := "tab_users"
	err = Dao.DeleteRow(&tab, idInt)
	if err != nil {
		fmt.Printf("删除用户信息 数据库 err:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"result": 0})
}

// 添加用户
func AddUsersHandler(c *gin.Context) {
	var newUser model.Users
	// name属性传值
	// id := c.PostForm("id")
	// username := c.PostForm("username")
	// password := c.PostForm("password")
	// fmt.Printf("id:%v,username:%v,password:%v\n", id, username, password)
	if err := c.ShouldBind(&newUser); err == nil { //{ID:1}
		fmt.Printf("newUser:%v\n", newUser)
		err := Dao.Adduser(&newUser)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// 成功则重定向到原页面
		c.Redirect(http.StatusMovedPermanently, "/admin/users")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func adminIsright(user string, passwd string) bool {
	var s string
	_ = Dao.DB.QueryRow("SELECT admin_id FROM tab_admin WHERE adminname =? and  admin_password =?", user, passwd).Scan(&s)
	if len(s) == 0 {
		return false
	} else {
		return true
	}
}
