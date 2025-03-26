package config

import (
	"os"
	"soundwave-go/internal/models"

	"gopkg.in/yaml.v3"
)

// UserConfig 用户配置
type UserConfig struct {
	Username    string              `yaml:"username"`
	Password    string              `yaml:"password"`
	Role        models.Role         `yaml:"role"`
	Permissions []models.Permission `yaml:"permissions"`
}

// MenuConfig 菜单配置
type MenuConfig struct {
	Name       string            `yaml:"name"`
	Path       string            `yaml:"path"`
	Icon       string            `yaml:"icon"`
	Permission models.Permission `yaml:"permission"`
	Sort       int               `yaml:"sort"`
}

// InitData 初始化数据结构
type InitData struct {
	Users []UserConfig `yaml:"users"`
	Menus []MenuConfig `yaml:"menus"`
}

func LoadInitData(path string) (*InitData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var initData InitData
	if err := yaml.Unmarshal(data, &initData); err != nil {
		return nil, err
	}

	return &initData, nil
}
