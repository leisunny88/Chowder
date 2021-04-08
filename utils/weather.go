package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func WeatherData(c *gin.Context, name string) {
	apiUrl := "http://apis.juhe.cn/simpleWeather/query"
	param := url.Values{}
	param.Set("city", name)
	param.Set("key", "e35e697f0ca7f863395aa020822cb263")

	// 发送请求
	data, err := Get(apiUrl, param)
	if err != nil {
		ResponseNotFoundCode(c, err.Error())
		return
	}
	var netReturn map[string]interface{}
	jsonErr := json.Unmarshal(data, &netReturn)
	if jsonErr != nil {
		ResponseNotFoundCode(c, jsonErr.Error())
		return
		//fmt.Errorf("请求异常:%v", jsonErr)
	} else {
		errorCode := netReturn["error_code"]
		reason := netReturn["reason"]
		data := netReturn["result"]
		allData := data.(map[string]interface{})["future"]
		fmt.Println(allData)
		// 当前天气信息
		realtime := data.(map[string]interface{})["realtime"]
		if errorCode.(float64) == 0 {
			// 请求成功，根据自身业务逻辑进行调整修改
			fmt.Printf("温度：%v\n湿度：%v\n天气：%v\n风向：%v\n风力：%v\n空气质量：%v",
				realtime.(map[string]interface{})["temperature"],
				realtime.(map[string]interface{})["humidity"],
				realtime.(map[string]interface{})["info"],
				realtime.(map[string]interface{})["direct"],
				realtime.(map[string]interface{})["power"],
				realtime.(map[string]interface{})["aqi"],
			)
		} else {
			// 查询失败，根据自身业务逻辑进行调整修改
			fmt.Printf("请求失败:%v_%v", errorCode.(float64), reason)
		}
	}

}

// get 方式发起网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
