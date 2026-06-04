package config

import (
	"os"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	UserCenter  UserCenterConfig  `yaml:"userCenter"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type UserCenterConfig struct {
	URL       string `yaml:"url"`
	SysName   string `yaml:"sysName"`
	JwtHashKey string `yaml:"jwtHashKey"`
}

var cfg *Config

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg = &Config{}
	if err := parseYaml(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}

func parseYaml(data []byte, v interface{}) error {
	// Simple YAML parser for basic config
	lines := splitLines(string(data))
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		parts := splitYamlLine(line)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "port":
			if cfg != nil {
				cfg.Server.Port = value
			}
		case "url":
			if cfg != nil {
				cfg.UserCenter.URL = value
			}
		case "sysName":
			if cfg != nil {
				cfg.UserCenter.SysName = value
			}
		case "jwtHashKey":
			if cfg != nil {
				cfg.UserCenter.JwtHashKey = value
			}
		}
	}
	return nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	indent := 0
	for i, c := range s {
		if c == '\n' {
			line := s[start:i]
			if indent == 0 || hasPrefix(line, "  ") {
				lines = append(lines, line)
			}
			start = i + 1
			indent = 0
		} else if c == ' ' && indent == 0 {
			indent = 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitYamlLine(line string) []string {
	idx := -1
	for i, c := range line {
		if c == ':' {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil
	}
	key := trimSpace(line[:idx])
	value := trimSpace(line[idx+1:])
	return []string{key, value}
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
