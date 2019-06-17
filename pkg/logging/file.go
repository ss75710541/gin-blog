package logging

import (
	"fmt"
	"gin-blog/pkg/setting"
	"time"
)

var (
	LogSavePath = setting.AppSetting.RuntimeRootPath
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}
