package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

// 密码回调函数
func (cfg *config) getPasswordCallback() func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	if !cfg.Auth.PasswordAuth.Enabled {
		return nil
	}
	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		log.Printf("New conncetion(src/dst): %s (%s)\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
		log.Printf("Authentication for ['%s','%s'] is accept\n", conn.User(), string(password))
		return nil, nil
	}
}

// ssh banner回调函数
func (cfg *config) getBannerCallback() func(conn ssh.ConnMetadata) string {
	if cfg.SSHProto.Banner == "" {
		return nil
	}
	banner := strings.ReplaceAll(strings.ReplaceAll(cfg.SSHProto.Banner, "\r\n", "\n"), "\n", "\r\n")
	if !strings.HasSuffix(banner, "\r\n") {
		banner = fmt.Sprintf("%v\r\n", banner)
	}
	return func(conn ssh.ConnMetadata) string { return banner }
}
