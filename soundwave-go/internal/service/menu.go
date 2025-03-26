package service

import (
	"context"
	"soundwave-go/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuService struct {
	db *mongo.Collection
}

func NewMenuService(db *mongo.Collection) *MenuService {
	return &MenuService{db: db}
}

func (s *MenuService) GetUserMenus(permissions []models.Permission) ([]models.Menu, error) {
	filter := bson.M{
		"permission": bson.M{
			"$in": permissions,
		},
	}

	cursor, err := s.db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var menus []models.Menu
	if err = cursor.All(context.Background(), &menus); err != nil {
		return nil, err
	}

	return menus, nil
}
