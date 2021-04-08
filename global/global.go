package global

import (
	"gin-project/config"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	CONFIG *config.Server  // 总配置信息
	VP     *viper.Viper
)
