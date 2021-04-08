package utils

import (
	"errors"
	"fmt"
	"gin-project/global"
)

// 模型数据创建
func CreateData(data interface{}) {
	global.DB.Create(data)
}

func UpdateData(data interface{}) {

}

// 模型数据查询
func SelectData(t interface{}, data interface{}) interface{} {
	result := global.DB.Find(&t)
	println(111112222)
	fmt.Println(result)
	return result
}

func BatchCreateData() {

}

// 模型数据删除
func DeleteData(m interface{}, id int) (err error) {
	global.DB.Where("id = ?", id).Take(&m)
	dbErr := global.DB.Delete(&m)
	if dbErr != nil {
		return errors.New("data deletion failed")
	}else {
		return nil
	}
}
