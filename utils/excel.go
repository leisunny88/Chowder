package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strconv"
)

func Excel() {
	router := gin.Default()
	// 给表单限制上传大小 (默认 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// 上传文件到指定的路径
		//c.SaveUploadedFile(file, "./"+file.Filename)
		xlFile, err := xlsx.OpenFile("./" + file.Filename)
		if err != nil {
			fmt.Println(err)
		}
		//var dataInfo []interface{}
		for _, sheet := range xlFile.Sheets {
			//遍历每一行
			for rowIndex, row := range sheet.Rows {
				//跳过第一行表头信息
				if rowIndex == 0 {
					// for _, cell := range row.Cells {
					//  text := cell.String()
					//  fmt.Printf("%s\n", text)
					// }
					continue
				}
				for _, cell := range row.Cells {
					text := cell.String()
					fmt.Printf("%s\n", text)
				}

			}
		}
	})

}

func UploadExcel(c *gin.Context) (err error, dataInfo []map[string]string) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	dst, _ := os.Getwd()
	dst = dst + "/static/excel/" + file.Filename
	fmt.Println(dst)
	// 上传文件到指定的路径
	_ = c.SaveUploadedFile(file, dst)
	xlFile, err := xlsx.OpenFile(dst)
	if err != nil {
		ResponseCodeMsg(c, "Incorrect format for uploading file")
		//fmt.Println(err)
	} else {
		//var dataInfo []map[string]string
		var DATA = make(map[string]string)

		for _, sheet := range xlFile.Sheets {
			//遍历每一行
			for rowIndex, row := range sheet.Rows {
				//跳过第一行表头信息
				if rowIndex == 0 {
					for i, cell := range row.Cells {
						text := cell.String()
						DATA[strconv.Itoa(i)] = text
						//fmt.Printf("%s\n", text)
					}
					continue
				}
				var dataLine = make(map[string]string)
				for j, cell := range row.Cells {
					key := DATA[strconv.Itoa(j)]
					text := cell.String()
					dataLine[key] = text
					//dataInfo = append(dataInfo, text)
				}
				if len(dataLine) != 0 {
					dataInfo = append(dataInfo, dataLine)
				}
			}
		}
	}
	rmErr := os.Remove(dst)
	if rmErr != nil {
		return errors.New("upload file cleaning failed"), nil
	} else {
		return nil, dataInfo
		//ResponseCodeMsg(c, "The uploaded file was cleaned up successfully")
	}
}
