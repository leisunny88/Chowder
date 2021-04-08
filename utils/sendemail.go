package utils

import (
	"gin-project/global"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

func SendEmail(c *gin.Context) {
	con := global.CONFIG.Email
	m := gomail.NewMessage()
	// 初始化发件人
	m.SetAddressHeader("From", con.From, con.Nickname)
	// 初始化收件人
	m.SetHeader("To",
		m.FormatAddress("xxxxx.@163.com", "xxx"),
		)
	// 主题
	m.SetHeader("Subject", "邀请函")
	// 正文
	m.SetBody("text/html", "hello")
	// 发送
	d := gomail.NewPlainDialer(con.Host, con.Port, con.From, con.Secret)
	if err := d.DialAndSend(m); err != nil {
		ResponseNotFoundCode(c, err.Error())
	}
}