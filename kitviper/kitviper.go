package kitviper

import (
    "github.com/spf13/viper"
)

// 读取本地toml配置文件并绑定到对象
func ReadToml(cfgFilePath string, bindToObj any) error {
    v := viper.New()
    v.SetConfigFile(cfgFilePath)
    v.SetConfigType("toml")
    err := v.ReadInConfig()
    if err != nil {
        return err
    }
    return v.Unmarshal(bindToObj)
}

// 读取远程consul上的toml配置文件并绑定到对象
func ReadRemoteTomlOnConsul(consulHost, consulKey string, bindToObj any) error {
    v := viper.New()
    err := v.AddRemoteProvider("consul", consulHost, consulKey)
    if err != nil {
        return err
    }
    v.SetConfigType("toml")
    err = v.ReadRemoteConfig()
    if err != nil {
        return err
    }
    return v.Unmarshal(bindToObj)
}
