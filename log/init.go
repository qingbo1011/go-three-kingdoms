package log

import (
	logging "github.com/sirupsen/logrus"
)

// Init 设置logrus的相关配置
func Init() {
	customFormatter := new(logging.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示
	customFormatter.ForceColors = true                      // 强制开启颜色

	logging.SetFormatter(customFormatter)
	//logging.SetOutput(os.Stdout)
	//logging.SetLevel(logging.DebugLevel)
}
