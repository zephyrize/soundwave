package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Menu 菜单模型
type Menu struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Path       string             `bson:"path" json:"path"`
	Icon       string             `bson:"icon" json:"icon"`
	Permission Permission         `bson:"permission" json:"permission"`
	Sort       int                `bson:"sort" json:"sort"`
}
