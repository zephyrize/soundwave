package service

import (
	"context"
	"errors"
	"soundwave-go/internal/models"
	"time"

	"soundwave-go/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	users *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		users: db.Collection("users"),
	}
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	logger.InfoLogger.Println("开始查询用户列表")

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := s.users.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.ErrorLogger.Printf("查询用户列表失败: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		logger.ErrorLogger.Printf("解析用户列表失败: %v", err)
		return nil, err
	}

	logger.InfoLogger.Printf("查询到 %d 个用户", len(users))
	return users, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// 检查用户名是否已存在
	count, err := s.users.CountDocuments(ctx, bson.M{"username": user.Username})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = s.users.InsertOne(ctx, user)
	return err
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id string, updates bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()
	_, err = s.users.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	return err
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 不允许删除管理员账户
	var user models.User
	err = s.users.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return err
	}
	if user.Role == models.RoleAdmin {
		return errors.New("不能删除管理员账户")
	}

	_, err = s.users.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.users.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ValidateUserInput 验证用户输入
func (s *UserService) ValidateUserInput(user *models.User) error {
	if user.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(user.Username) < 3 || len(user.Username) > 32 {
		return errors.New("用户名长度必须在3-32个字符之间")
	}
	if user.Password != "" && (len(user.Password) < 6 || len(user.Password) > 32) {
		return errors.New("密码长度必须在6-32个字符之间")
	}
	if user.Role == "" {
		return errors.New("用户角色不能为空")
	}
	if !isValidRole(user.Role) {
		return errors.New("无效的用户角色")
	}
	return nil
}

// isValidRole 检查角色是否有效
func isValidRole(role models.Role) bool {
	validRoles := []models.Role{models.RoleAdmin, models.RoleUser, models.RoleTester}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// UpdateUserPassword 更新用户密码
func (s *UserService) UpdateUserPassword(ctx context.Context, id string, newPassword string) error {
	if len(newPassword) < 6 || len(newPassword) > 32 {
		return errors.New("密码长度必须在6-32个字符之间")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now(),
		},
	}

	_, err = s.users.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger.InfoLogger.Printf("查询用户名: %s", username)

	var user models.User
	err := s.users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.InfoLogger.Printf("用户名 %s 不存在", username)
			return nil, nil
		}
		logger.ErrorLogger.Printf("查询用户失败: %v", err)
		return nil, err
	}

	logger.InfoLogger.Printf("找到用户: %s", username)
	return &user, nil
}
