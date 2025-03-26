package service

import (
	"context"
	"errors"
	"fmt"
	"soundwave-go/internal/models"
	"time"

	"soundwave-go/internal/config"
	"soundwave-go/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users  *mongo.Collection
	tokens *mongo.Collection
	config *config.Config
}

func NewAuthService(db *mongo.Database, cfg *config.Config) *AuthService {
	return &AuthService{
		users:  db.Collection(cfg.MongoDB.Collections.Users),
		tokens: db.Collection(cfg.MongoDB.Collections.Tokens),
		config: cfg,
	}
}

func (s *AuthService) Register(user *models.User) error {
	// 检查用户名是否已存在
	exists, err := s.users.CountDocuments(context.Background(), bson.M{"username": user.Username})
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 设置默认角色和权限
	user.Role = models.RoleUser
	user.Permissions = []models.Permission{models.PermissionViewServices, models.PermissionViewStats}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = s.users.InsertOne(context.Background(), user)
	return err
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	var user models.User
	err := s.users.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	return utils.GenerateToken(s.config, user)
}

func (s *AuthService) ChangePassword(userID string, oldPassword, newPassword string) error {
	var user models.User
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("无效的用户ID")
	}

	// 查找用户
	err = s.users.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 生成新密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	update := bson.M{
		"$set": bson.M{
			"password":  string(hashedPassword),
			"updatedAt": time.Now(),
		},
	}

	_, err = s.users.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	return err
}
