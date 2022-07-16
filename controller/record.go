package controller

import (
	"fmt"
	"golibrary/Dao"
	"golibrary/gettime"
	"golibrary/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// record 页面
func RecordsHander(c *gin.Context) {
	//claims := c.MustGet("claims").(*md.CustomClaims)
	//if claims != nil {
	//	fmt.Println("测试recod页面")
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": 0,
	//		"msg":    "token有效",
	//		"data":   claims,
	//	})
	//}
	c.HTML(http.StatusOK, "admin/record.html", gin.H{
		"num":     1,
		"data":    nil,
		"nowTime": gettime.GetTime(),
	})
	//recordAll, err := Dao.SelectAllrecord()
	//if err != nil {
	//	panic("tab_record select err")
	//}
	//// fmt.Printf("books:%#v\n", booksList) //切片数组
	//// for i, v := range booksList {
	//// 	fmt.Println(i, v)
	//// }
	//if len(*recordAll) > 0 {
	//	c.HTML(http.StatusOK, "admin/record.html", gin.H{
	//		"num":     1,
	//		"data":    recordAll,
	//		"nowTime": getTime.GetTime(),
	//	})
	//} else {
	//	c.HTML(http.StatusOK, "admin/record.html", gin.H{
	//		"num":     1,
	//		"data":    nil,
	//		"nowTime": getTime.GetTime(),
	//	})
	//}

}

// 添加一条记录
func AddRecordHandler(c *gin.Context) {
	var newRecord model.Record
	if err := c.ShouldBind(&newRecord); err == nil { //{ID:1}
		fmt.Printf("newUser:%v\n", newRecord)
		err := Dao.AddRecordDB(&newRecord)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// 成功则重定向到原页面
		c.Redirect(http.StatusMovedPermanently, "/admin/record")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// 删除一条借阅记录
func DeleteRecordHandler(c *gin.Context) {
	idStr := c.PostForm("id")
	// string到int
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("删除借阅记录 err:%v\n", err)
	}
	// 删除借阅记录
	tab := "tab_record"
	err = Dao.DeleteRow(&tab, idInt)
	if err != nil {
		fmt.Printf("删除借阅记录 数据库 err:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"result": 0})
}
