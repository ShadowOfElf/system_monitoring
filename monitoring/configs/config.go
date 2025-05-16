package configs

import (
	"fmt"
	"net"

	"github.com/ShadowOfElf/system_monitoring/internal/logger"
	"github.com/ShadowOfElf/system_monitoring/internal/resources"
	"github.com/spf13/viper"
)

type GRPCConf struct {
	Addr string
}

type LoggerConf struct {
	Level logger.LogLevel
}

type Config struct {
	GRPC          GRPCConf
	Logger        LoggerConf
	Enable        resources.CollectorEnable
	MaxSize       int
	RepeatRateSec int
}

func NewConfig(configFile string) *Config {
	var err error
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Configuration is not loaded, default values will be used")
	}

	maxSize := viper.GetInt("app.max_size")
	repeatRate := viper.GetInt("app.repeat_rate")

	logLevel := viper.GetString("logger.level")
	logValid := logLevelValidator(logLevel)
	if logLevel == "" || !logValid {
		fmt.Println("log automatic set to default")
		logLevel = "INFO"
	}

	grpcHost := viper.GetString("grpc.host")
	grpcPort := viper.GetString("grpc.port")

	addrGRPC := net.JoinHostPort(grpcHost, grpcPort)
	_, err = net.ResolveTCPAddr("tcp", addrGRPC)
	if err != nil {
		fmt.Println("host or port GRPC incorrect, using default")
		addrGRPC = "127.0.0.1:8070"
	}

	enable := resources.CollectorEnable{
		Load:       viper.GetBool("collector.load"),
		CPU:        viper.GetBool("collector.cpu"),
		Disk:       viper.GetBool("collector.disk"),
		Net:        viper.GetBool("collector.net"),
		TopTalkers: viper.GetBool("collector.top_talkers"),
	}

	return &Config{
		Logger:        LoggerConf{Level: logger.LogLevel(logLevel)},
		GRPC:          GRPCConf{Addr: addrGRPC},
		Enable:        enable,
		MaxSize:       maxSize,
		RepeatRateSec: repeatRate,
	}
}

func logLevelValidator(level string) bool {
	allowLevel := map[string]bool{
		string(logger.DebugLevel): true,
		string(logger.WarnLevel):  true,
		string(logger.InfoLevel):  true,
		string(logger.ErrorLevel): true,
	}
	return allowLevel[level]
}
