package main

import (
	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"

	"flag"
	"log"
	"net"
	"path"
)

func main() {
	// hostkey文件所在路径
	dataDir := flag.String("data_dir", path.Join(xdg.DataHome, "hostkeys"), "data directory")
	flag.Parse()
	configString := ""

	// 获取ssh连接的配置文件
	cfg, err := getConfig(configString, *dataDir)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	// 监听端口
	listener, err := net.Listen("tcp", cfg.Server.ListenAddress)
	if err != nil {
		log.Fatalf("Failed to listen for connections: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %v", listener.Addr())

	// 接收所有请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		// 设置连接操作
		go handleConnection(conn, cfg)
	}
}

// 获取配置文件
func getConfig(configString string, dataDir string) (*config, error) {
	// 1.获取默认配置文件
	cfg := getDefaultConfig()

	if err := yaml.UnmarshalStrict([]byte(configString), cfg); err != nil {
		return nil, err
	}

	// 2.判断主机密钥是否为空  如果为空这设置默认主机密钥
	if len(cfg.Server.HostKeys) == 0 {
		log.Printf("No host keys configured, using keys at %q", dataDir)
		if err := cfg.setDefaultHostKeys(dataDir); err != nil {
			return nil, err
		}
	}

	// 3.设置ssh配置文件
	if err := cfg.setupSSHConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
