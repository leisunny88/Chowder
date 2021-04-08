package initialize

import (
	"fmt"
	"gin-project/global"
	"gin-project/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Migrate(c *gin.Context) {
	err := global.DB.AutoMigrate(
		model.Device{},
		model.User{},
		model.Disk{},
		model.Memory{},
		model.CPU{},
	)
	if err != nil {
		panic(err)
	}
	//os.Exit(0)
}

//type MysqlConnectPool struct {
//}
//
//func (m *MysqlConnectPool) InitDataPool() (issucc bool) {
//	dsn := "Gin-demo:iXaHBxj04cydLoOB@(192.168.100.120:3306)/Gin-demo?charset=utf8mb4&parseTime=True&loc=Local"
//	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//		return false
//	}
//	//关闭数据库，db会被多个goroutine共享，可以不调用
//	//defer db.Close()
//	return true
//}
//
//func (m *MysqlConnectPool) GetMysqlDB() (db *gorm.DB) {
//	return db
//}

// 数据模型初始化
func InitDB() (db *gorm.DB) {
	m := global.CONFIG.Mysql
	//dsn := "root:leisunny@(127.0.0.1:3306)/Gin_test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("wwww")
	}
	return db
}
