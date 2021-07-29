package main

import (
	"golang.org/x/crypto/ssh"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// server 配置文件 对应yaml文件中的server
type serverConfig struct {
	ListenAddress string   `yaml:"listen_address"`
	HostKeys      []string `yaml:"host_keys"`
}

// 日志配置文件 对应yaml文件中的logging
type loggingConfig struct {
	File       string `yaml:"file"`
	JSON       bool   `yaml:"json"`
	Timestamps bool   `yaml:"timestamps"`
	Debug      bool   `yaml:"debug"`
}

// 认证配置文件 对应yaml文件中的auth
type authConfig struct {
	MaxTries     int              `yaml:"max_tries"`
	NoAuth       bool             `yaml:"no_auth"`
	PasswordAuth commonAuthConfig `yaml:"password_auth"`
	//PublicKeyAuth commonAuthConfig `yaml:"public_key_auth"`
}

// 认证的两种情况
type commonAuthConfig struct {
	Enabled  bool `yaml:"enabled"`
	Accepted bool `yaml:"accepted"`
}

// ssh 协议的配置文件 对应yaml文件中的ssh_proto
type sshProtoConfig struct {
	Version string `yaml:"version"`
	Banner  string `yaml:"banner"`
}

// 整个ssh的配置
type config struct {
	Server   serverConfig   `yaml:"server"`
	Logging  loggingConfig  `yaml:"logging"`
	Auth     authConfig     `yaml:"auth"`
	SSHProto sshProtoConfig `yaml:"ssh_proto"`

	parsedHostKeys []ssh.Signer // 存放解析后的主机密钥
	sshConfig      *ssh.ServerConfig
}

// 1.默认配置文件
func getDefaultConfig() *config {
	cfg := &config{}
	cfg.Server.ListenAddress = "127.0.0.1:2222"
	cfg.Logging.Timestamps = true
	cfg.Auth.PasswordAuth.Enabled = true
	cfg.Auth.PasswordAuth.Accepted = true
	//cfg.Auth.PublicKeyAuth.Enabled = true
	cfg.SSHProto.Version = "SSH-2.0-gossh-honey"
	cfg.SSHProto.Banner = "This is an SSH honeypot. Everything is logged and monitored."
	return cfg
}

// 2.设置默认的主机密钥
func (cfg *config) setDefaultHostKeys(dataDir string) error {
	keyFile, err := generateKey(dataDir) // 在指定路径中生成hostkey
	if err != nil {
		return nil
	}
	cfg.Server.HostKeys = append(cfg.Server.HostKeys, keyFile)
	return nil
}

// 2.1生成密钥
func generateKey(dataDir string) (string, error) {
	keyFile := path.Join(dataDir, "host_rsa_key")
	if _, err := os.Stat(keyFile); err == nil {
		return keyFile, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	log.Printf("Host key %q not found, generating it", keyFile)

	if _, err := os.Stat(path.Dir(keyFile)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(keyFile), 0755); err != nil {
			return "", err
		}
	}

	var key interface{}

	// 调用rsa的GenerateKey方法产生密钥
	key, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return "", err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}), 0600); err != nil {
		return "", err
	}
	return keyFile, nil
}

// 3.获取ssh 配置文件
func (cfg *config) setupSSHConfig() error {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth:     cfg.Auth.NoAuth,
		MaxAuthTries:     cfg.Auth.MaxTries,
		PasswordCallback: cfg.getPasswordCallback(),
		//PublicKeyCallback: cfg.getPublicKeyCallback(),
		BannerCallback: cfg.getBannerCallback(),
		ServerVersion:  cfg.SSHProto.Version,
	}
	// 3.1解析主机密钥
	if err := cfg.parseHostKeys(); err != nil {
		return err
	}
	for _, key := range cfg.parsedHostKeys {
		// 添加主机密钥
		sshConfig.AddHostKey(key)
	}
	cfg.sshConfig = sshConfig
	return nil
}

// 3.1解析密钥
func (cfg *config) parseHostKeys() error {
	for _, keyFile := range cfg.Server.HostKeys {
		signer, err := loadKey(keyFile) // 加载keyFile
		if err != nil {
			return err
		}
		cfg.parsedHostKeys = append(cfg.parsedHostKeys, signer)
	}
	return nil
}

// 3.2加载密钥
func loadKey(keyFile string) (ssh.Signer, error) {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(keyBytes)
}
