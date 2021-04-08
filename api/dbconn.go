package api

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	dsn := "Gin-demo:iXaHBxj04cydLoOB@(192.168.100.120:3306)/Gin-demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("wwww")
	}
	//initialize.MysqlTables(db)
	//c.JSON(200, gin.H{
	//	"msg": "init db success",
	//})
	return db
}
