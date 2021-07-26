package main

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"net"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
)

// ConfigInstance 对应配置文件的结构体
type ConfigInstance struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Command  string `json:"command"`
	Port     int    `json:"port"`
}

var configInstance ConfigInstance

func main() {
	// 配置文件初始化
	InitConfig()

	// SSH 服务端
	config := ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// 验证用户名和密码是否正确
			if conn.User() == configInstance.Name && string(password) == configInstance.Password {
				log.Println("Login successful. username:", conn.User(), "address:", conn.RemoteAddr().String())
				return nil, nil
			} else {
				log.Println("Login failed. username:", conn.User(), "password:", string(password), "address:", conn.RemoteAddr().String())
				return nil, errors.New("Unknow username or password")
			}
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key:", err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	// 添加主机密钥
	config.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(configInstance.Port))
	if err != nil {
		log.Fatal("failed to listen for connection:", err)
	}
	nConn, _ := listener.Accept()

	_, _, _, err = ssh.NewServerConn(nConn, &config)
	if err != nil {
		fmt.Printf("error is: %v", err)
	}

}

// InitConfig 初始化配置的操作
func InitConfig() {
	// 根据路径读取配置文件信息
	if contents, err := ioutil.ReadFile(filepath.Join("config", "config.json")); err == nil {
		log.Println("Load config:", filepath.Join("config", "config.json"))

		err := json.Unmarshal(contents, &configInstance)
		if err != nil {
			log.Println("Load config error:", err)
			panic(err)
		}
	} else {
		if err != nil {
			log.Println("Load config error:")
			panic(err)
		}
	}
	log.Println("[Config] Username:", configInstance.Name)
	log.Println("[Config] Port:", configInstance.Port)
}
