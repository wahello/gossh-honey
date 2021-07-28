package main

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
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
	MaxTries      int              `yaml:"max_tries"`
	NoAuth        bool             `yaml:"no_auth"`
	PasswordAuth  commonAuthConfig `yaml:"password_auth"`
	PublicKeyAuth commonAuthConfig `yaml:"public_key_auth"`
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

	parsedHostKeys []ssh.Signer
	sshConfig      *ssh.ServerConfig
	logFileHandle  io.WriteCloser
}

// 默认配置文件
func getDefaultConfig() *config {
	cfg := &config{}
	cfg.Server.ListenAddress = "127.0.0.1:2222"
	cfg.Logging.Timestamps = true
	cfg.Auth.PasswordAuth.Enabled = true
	cfg.Auth.PasswordAuth.Accepted = true
	cfg.Auth.PublicKeyAuth.Enabled = true
	cfg.SSHProto.Version = "SSH-2.0-gossh-honey"
	cfg.SSHProto.Banner = "This is an SSH honeypot. Everything is logged and monitored."
	return cfg
}

// 生成密钥
func generateKey(dataDir string) (string, error) {
	keyFile := path.Join(dataDir, "host_rsa_key")
	if _, err := os.Stat(keyFile); err == nil {
		return keyFile, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}
	infoLogger.Printf("Host key %q not found, generating it", keyFile)
	if _, err := os.Stat(path.Dir(keyFile)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(keyFile), 0755); err != nil {
			return "", err
		}
	}
	var key interface{}
	err := errors.New("unsupported key type")
	// 调用rsa的GenerateKey方法产生密钥
	key, err = rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return "", err
	}

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

// 加载密钥
func loadKey(keyFile string) (ssh.Signer, error) {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(keyBytes)
}

// 设置默认的主机密钥
func (cfg *config) setDefaultHostKeys(dataDir string) error {
	keyFile, err := generateKey(dataDir)
	if err != nil {
		return nil
	}
	cfg.Server.HostKeys = append(cfg.Server.HostKeys, keyFile)
	return nil
}

func (cfg *config) parseHostKeys() error {
	for _, keyFile := range cfg.Server.HostKeys {
		signer, err := loadKey(keyFile)
		if err != nil {
			return err
		}
		cfg.parsedHostKeys = append(cfg.parsedHostKeys, signer)
	}
	return nil
}

// 获取ssh 配置文件
func (cfg *config) setupSSHConfig() error {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth:      cfg.Auth.NoAuth,
		MaxAuthTries:      cfg.Auth.MaxTries,
		PasswordCallback:  cfg.getPasswordCallback(),
		PublicKeyCallback: cfg.getPublicKeyCallback(),
		ServerVersion:     cfg.SSHProto.Version,
		BannerCallback:    cfg.getBannerCallback(),
	}
	if err := cfg.parseHostKeys(); err != nil {
		return err
	}
	for _, key := range cfg.parsedHostKeys {
		sshConfig.AddHostKey(key)
	}
	cfg.sshConfig = sshConfig
	return nil
}

// 设置日志文件
func (cfg *config) setupLogging() error {
	if cfg.logFileHandle != nil {
		cfg.logFileHandle.Close()
	}
	if cfg.Logging.File != "" {
		logFile, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
		cfg.logFileHandle = logFile
	} else {
		log.SetOutput(os.Stdout)
		cfg.logFileHandle = nil
	}
	if !cfg.Logging.JSON && cfg.Logging.Timestamps {
		log.SetFlags(log.LstdFlags)
	} else {
		log.SetFlags(0)
	}
	return nil
}

// 获取配置文件
func getConfig(configString string, dataDir string) (*config, error) {
	cfg := getDefaultConfig()

	if err := yaml.UnmarshalStrict([]byte(configString), cfg); err != nil {
		return nil, err
	}

	if len(cfg.Server.HostKeys) == 0 {
		infoLogger.Printf("No host keys configured, using keys at %q", dataDir)
		// 设置默认密钥
		if err := cfg.setDefaultHostKeys(dataDir); err != nil {
			return nil, err
		}
	}

	// 设置配置文件
	if err := cfg.setupSSHConfig(); err != nil {
		return nil, err
	}

	// 设置日志文件
	if err := cfg.setupLogging(); err != nil {
		return nil, err
	}

	return cfg, nil
}
