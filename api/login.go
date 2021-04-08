package api

import (
	"errors"
	"gin-project/global"
	"gin-project/model"
	"gin-project/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("pwd")
	users := model.User{}
	// 数据库查询错误时的处理
	err := global.DB.Where("name = ?", username).Take(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ResponseCodeMsg(c, "The user name does not exist")
	}
	//fmt.Println(reflect.TypeOf(users))   打印类型
	OrmUser := users.Name
	OrmPwd := users.PassWord
	if len(OrmUser) != 0 {
		err := bcrypt.CompareHashAndPassword([]byte(OrmPwd), []byte(pwd)) // 密码校验
		if err != nil {
			utils.ResponseCodeMsg(c, "Login password error")
		} else {
			utils.ResponseCodeMsg(c, "Login successful")
		}
	}
}
