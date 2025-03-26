package db

import (
	"context"

	"soundwave-go/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// SaveService 保存服务实例到MongoDB
func (m *MongoDB) SaveService(key string, service interface{}) error {
	coll := m.Collection("services")

	// 使用uniqueID作为文档ID
	filter := bson.M{"_id": key}
	update := bson.M{"$set": service}
	_, err := coll.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		logger.ErrorLogger.Printf("保存服务到MongoDB失败：%v", err)
		return err
	}

	return nil
}

// GetService 从MongoDB获取服务实例
func (m *MongoDB) GetService(key string, result interface{}) error {
	coll := m.Collection("services")

	err := coll.FindOne(context.Background(), bson.M{"_id": key}).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}

	return nil
}

// DeleteService 从MongoDB删除服务实例
func (m *MongoDB) DeleteService(key string) error {
	coll := m.Collection("services")

	_, err := coll.DeleteOne(context.Background(), bson.M{"_id": key})
	return err
}

// ListServices 获取所有服务实例
func (m *MongoDB) ListServices(result interface{}) error {
	coll := m.Collection("services")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	return cursor.All(context.Background(), result)
}

func NewMongoDB(uri, database string) (*MongoDB, error) {
	logger.InfoLogger.Printf("连接 MongoDB: %s, 数据库: %s", uri, database)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.ErrorLogger.Printf("MongoDB 连接失败: %v", err)
		return nil, err
	}

	// 测试连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.ErrorLogger.Printf("MongoDB Ping 失败: %v", err)
		return nil, err
	}

	logger.InfoLogger.Println("MongoDB 连接成功")
	return &MongoDB{
		client: client,
		db:     client.Database(database),
	}, nil
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

func (m *MongoDB) Database() *mongo.Database {
	return m.db
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}
