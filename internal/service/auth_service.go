package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/cycling02/ai-novel-backend/internal/model"
	"github.com/cycling02/ai-novel-backend/internal/repository"
)

var (
	ErrUserExists      = errors.New("用户已存在")
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrInvalidToken    = errors.New("无效的令牌")
	ErrTokenExpired    = errors.New("令牌已过期")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtSecret  []byte
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, string, error) {
	// 检查用户是否已存在
	existing, _ := s.userRepo.GetByEmail(req.Email)
	if existing != nil {
		return nil, "", ErrUserExists
	}

	// 创建用户
	user := &model.User{
		Email:    req.Email,
		Nickname: req.Nickname,
		Role:     "user",
		Status:   "active",
	}

	// 加密密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, "", err
	}

	// 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// 生成 Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", ErrUserNotFound
	}

	// 检查密码
	if !user.CheckPassword(req.Password) {
		return nil, "", ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status == "banned" {
		return nil, "", errors.New("账户已被禁用")
	}

	// 生成 Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	// 解析 Token
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", ErrUserNotFound
	}

	// 生成新 Token
	return s.generateToken(user)
}

// ValidateToken 验证 Token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*model.AuthClaims, error) {
	return s.parseToken(tokenString)
}

// generateToken 生成 JWT Token
func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := model.AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 天过期
		"iat":     time.Now().Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

// parseToken 解析 JWT Token
func (s *AuthService) parseToken(tokenString string) (*model.AuthClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 检查是否过期
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, ErrTokenExpired
			}
		}

		return &model.AuthClaims{
			UserID: claims["user_id"].(string),
			Email:  claims["email"].(string),
			Role:   claims["role"].(string),
		}, nil
	}

	return nil, ErrInvalidToken
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}