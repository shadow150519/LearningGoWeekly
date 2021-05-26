package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("userInfo")
	viper.AddConfigPath("./demo/config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(viper.Get("User1.Username"))
	// 写入到一个新文件中
	viper.Set("User1.Age",30)
	viper.SafeWriteConfigAs("./demo/config/newUserInfo.yaml")

	// 要在WatchConfig前调用这个方法
	// 一次改变会调用两次这个函数，这是个viper自己的bug
	viper.OnConfigChange(func(in fsnotify.Event){
		fmt.Println("配置改变了")
	})
	viper.WatchConfig()
	for  {

	}
	return
}
