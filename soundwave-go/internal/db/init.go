package db

import (
	"context"
	"time"

	"soundwave-go/internal/config"
	"soundwave-go/internal/logger"
	"soundwave-go/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitializeData(db *mongo.Database, initData *config.InitData) error {
	logger.InfoLogger.Println("开始初始化数据...")

	// 清空现有数据
	if err := clearCollections(db); err != nil {
		logger.ErrorLogger.Printf("清空集合失败: %v", err)
		return err
	}
	logger.InfoLogger.Println("清空集合完成")

	// 初始化用户数据
	if err := initUsers(db, initData.Users); err != nil {
		logger.ErrorLogger.Printf("初始化用户数据失败: %v", err)
		return err
	}
	logger.InfoLogger.Println("初始化用户数据完成")

	// 初始化菜单数据
	if err := initMenus(db, initData.Menus); err != nil {
		logger.ErrorLogger.Printf("初始化菜单数据失败: %v", err)
		return err
	}
	logger.InfoLogger.Println("初始化菜单数据完成")

	return nil
}

func clearCollections(db *mongo.Database) error {
	collections := []string{"users", "menus"}
	for _, name := range collections {
		if err := db.Collection(name).Drop(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func initUsers(db *mongo.Database, users []config.UserConfig) error {
	if len(users) == 0 {
		logger.WarnLogger.Println("没有用户数据需要初始化")
		return nil
	}

	var documents []interface{}
	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		doc := models.User{
			ID:          primitive.NewObjectID(),
			Username:    user.Username,
			Password:    string(hashedPassword),
			Role:        user.Role,
			Permissions: user.Permissions,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		documents = append(documents, doc)
		logger.InfoLogger.Printf("准备创建用户: %s, 角色: %s", user.Username, user.Role)
	}

	result, err := db.Collection("users").InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}
	logger.InfoLogger.Printf("成功创建 %d 个用户", len(result.InsertedIDs))
	return nil
}

func initMenus(db *mongo.Database, menus []config.MenuConfig) error {
	if len(menus) == 0 {
		return nil
	}

	var documents []interface{}
	for _, menu := range menus {
		documents = append(documents, models.Menu{
			Name:       menu.Name,
			Path:       menu.Path,
			Icon:       menu.Icon,
			Permission: menu.Permission,
			Sort:       menu.Sort,
		})
	}

	_, err := db.Collection("menus").InsertMany(context.Background(), documents)
	return err
}
