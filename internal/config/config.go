package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Vector   VectorConfig   `mapstructure:"vector"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type VectorConfig struct {
	Provider  string `mapstructure:"provider"`
	APIKey    string `mapstructure:"api_key"`
	IndexName string `mapstructure:"index_name"`
	Namespace string `mapstructure:"namespace"`
}

type LLMConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"`
	Model    string `mapstructure:"model"`
}

type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	TokenExpiry   string `mapstructure:"token_expiry"`
	RefreshExpiry string `mapstructure:"refresh_expiry"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用环境变量
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖
	if url := os.Getenv("DATABASE_URL"); url != "" {
		cfg.Database.URL = url
	}
	if key := os.Getenv("PINECONE_API_KEY"); key != "" {
		cfg.Vector.APIKey = key
	}
	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		cfg.LLM.APIKey = key
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		cfg.LLM.APIKey = key
		cfg.LLM.Provider = "openai"
	}

	return &cfg, nil
}
