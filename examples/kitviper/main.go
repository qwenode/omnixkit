package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/qwenode/omnixkit/kitviper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

type DatabaseConfig struct {
	DSN          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

func main() {
	configPath := flag.String("config", "examples/kitviper/config.toml", "toml 配置文件路径")
	flag.Parse()

	var cfg Config
	if err := kitviper.ReadToml(*configPath, &cfg); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	addr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port))
	fmt.Println("server.addr:", addr)
	fmt.Println("server.debug:", cfg.Server.Debug)
	fmt.Println("database.dsn:", cfg.Database.DSN)
	fmt.Println("database.max_open_conns:", cfg.Database.MaxOpenConns)
	fmt.Println("database.max_idle_conns:", cfg.Database.MaxIdleConns)
}

