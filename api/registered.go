package api

import (
	"gin-project/global"
	"gin-project/model"
	"gin-project/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 注册用户功能
func Register(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("pwd")
	// 加密处理
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	encodePWD := string(hash)
	if len(username) == 0 && len(pwd) == 0 {
		utils.ResponseCodeMsg(c, "username and password is null")
	} else if len(username) < 6 && len(pwd) < 8 {
		utils.ResponseCodeMsg(c, "The user name is less than 6 characters long and "+
			"the password is less than 8 characters long")
	} else {
		user := model.User{Name: username, PassWord: encodePWD}
		//utils.CreateData(&user)}
		//resultArr := make(map[interface{}]Result)
		//global.DB.Table("users").Select([]string{"Name", "pass_word"}).Scan(&resultArr)
		var users []model.User
		global.DB.Select("name, pass_word").Find(&users)
		//fmt.Println(users)
		if len(users) == 0 {
			user.Name = username
			user.PassWord = pwd
			println(111111)
			utils.CreateData(&user)
		} else {
			for k, v := range users {
				println(k, v.Name)
				userName := v.Name
				if username == userName {
					utils.ResponseCodeMsg(c, "The user name already exists")
				} else {
					utils.CreateData(&user)
				}
			}
		}
	}
}

type Result struct {
	Name     string
	PassWord string
}
