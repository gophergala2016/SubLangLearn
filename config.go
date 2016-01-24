package main

import (
	"os"
	"io/ioutil"
	"strings"
	"strconv"
)

type Config struct {
	VlcPlayerPath string
	VlcPlayerTcpPort int
	WebServerHttpPort int
}

func LoadConfig(path string) *Config {
	config := &Config{VlcPlayerPath:`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`, VlcPlayerTcpPort:3016, WebServerHttpPort:2016}
	if _, err := os.Stat(path); err == nil {
		content, err := ioutil.ReadFile(path)
		if err == nil {
			lines := strings.Split(string(content),"\n")
			for _, line := range lines {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) < 2 {
					continue
				}
				key := strings.ToLower(strings.TrimSpace(parts[0]))
				value := strings.ToLower(strings.TrimSpace(parts[1]))
				if key == "vlc_path" && value != "" {
					config.VlcPlayerPath = value
				} else if key == "vlc_port" && value != "" {
					port, _ := strconv.Atoi(value)
					if port != 0 {
						config.VlcPlayerTcpPort = port
					}
				} else if key == "web_port" && value != "" {
					port, _ := strconv.Atoi(value)
					if port != 0 {
						config.WebServerHttpPort = port
					}
				}
			}
		}
	}
	return config
}
