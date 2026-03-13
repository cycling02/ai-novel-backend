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
	Provider string `mapstructure:"provider"` // deepseek, openai, minimax, zhipu, moonshot
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

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用环境变量
	}

	// 环境变量优先
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

	// LLM 配置支持多个提供商
	if provider := os.Getenv("LLM_PROVIDER"); provider != "" {
		cfg.LLM.Provider = provider
	}
	if baseURL := os.Getenv("LLM_BASE_URL"); baseURL != "" {
		cfg.LLM.BaseURL = baseURL
	}
	if model := os.Getenv("LLM_MODEL"); model != "" {
		cfg.LLM.Model = model
	}

	// 支持不同提供商的 API Key
	if cfg.LLM.Provider == "minimax" {
		if key := os.Getenv("MINIMAX_API_KEY"); key != "" {
			cfg.LLM.APIKey = key
		}
	} else if cfg.LLM.Provider == "zhipu" {
		if key := os.Getenv("ZHIPU_API_KEY"); key != "" {
			cfg.LLM.APIKey = key
		}
	} else if cfg.LLM.Provider == "moonshot" {
		if key := os.Getenv("MOONSHOT_API_KEY"); key != "" {
			cfg.LLM.APIKey = key
		}
	} else if cfg.LLM.Provider == "openai" {
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			cfg.LLM.APIKey = key
		}
	} else {
		// 默认使用 DeepSeek
		if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
			cfg.LLM.APIKey = key
		}
	}

	return &cfg, nil
}

// GetDefaultLLMConfig 获取默认 LLM 配置
func GetDefaultLLMConfig(provider string) LLMConfig {
	configs := map[string]LLMConfig{
		"deepseek": {
			Provider: "deepseek",
			BaseURL:  "https://api.deepseek.com",
			Model:    "deepseek-chat",
		},
		"openai": {
			Provider: "openai",
			BaseURL:  "https://api.openai.com/v1",
			Model:    "gpt-4o",
		},
		"minimax": {
			Provider: "minimax",
			BaseURL:  "https://api.minimaxi.com/v1",
			Model:    "MiniMax-M2.5",
		},
		"zhipu": {
			Provider: "zhipu",
			BaseURL:  "https://open.bigmodel.cn/api/paas/v4",
			Model:    "glm-5",
		},
		"moonshot": {
			Provider: "moonshot",
			BaseURL:  "https://api.moonshot.cn/v1",
			Model:    "kimi-k2.5",
		},
	}

	if cfg, ok := configs[provider]; ok {
		return cfg
	}
	return configs["deepseek"]
}
