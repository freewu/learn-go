package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type User struct {
	ID   		int64		`` // 自增ID
	Age 		uint		`` // 年龄
	Name 		string		`` // 姓名
	CreatedAt 	time.Time 	`` // 创建时间
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:123456@192.168.150.129:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("connect error: ",err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("ping error: ",err)
	}
	fmt.Println("Connected to MongoDB!")

	// 连接库&表
	collection := client.Database("test").Collection("user")

	user := &User{ ID: 1, Age: 18,Name: "bluefrog", CreatedAt: time.Now() }
	// insert
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal("insert error: ",err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID) // ObjectID("61501d2139f294ae7cb53182")
	fmt.Printf("Inserted a single document: %#v \n", insertResult) // &{ObjectID("61501d2139f294ae7cb53182")}

	// batch insert
	batchUser := []interface{}{
		&User{ ID: 1, Age: 118, Name: "bluefrog1", CreatedAt: time.Now() },
		&User{ ID: 1, Age: 128, Name: "bluefrog2", CreatedAt: time.Now() },
	}
	insertManyResult, err := collection.InsertMany(context.TODO(),batchUser)
	if err != nil {
		log.Fatal("batch insert error: ",err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) // [ObjectID("61501d2139f294ae7cb53183") ObjectID("61501d2139f294ae7cb53184")]

	// update
	filter := bson.D{{"name", "bluefrog"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	// update := bson.M{"$set": server.Task_info{Task_id: "update task id"}}  // 不推荐直接用结构体，玩意结构体字段多了，初始化为零值。
	// 因为可能会吧零值更新到数据库，而不是像 gorm 的updates 忽略零值
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal("update error: ",err)
	}
	fmt.Println("Update documents: ",updateResult) // &{1 1 0 <nil>}
	// MatchedCount: 1 ModifiedCount: 1 UpsertedCount: 0 UpsertedID: <nil>
	fmt.Printf("MatchedCount: %v ModifiedCount: %v UpsertedCount: %v UpsertedID: %v\n",updateResult.MatchedCount,updateResult.ModifiedCount,updateResult.UpsertedCount,updateResult.UpsertedID)

	// 查询单个文档
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal("FindOne error: ",err)
	}
	fmt.Printf("Found a single document: %+v\n", result) // Found a single document: {ID:1 Age:20 Name:bluefrog CreatedAt:2021-09-26 06:36:34.245 +0000 UTC}

	// 查询多个文档
	one, _:= collection.Find(context.TODO(), bson.M{"name": "bluefrog"})
	defer func() {  // 关闭
		if err := one.Close(context.TODO()); err != nil {
			log.Fatal("close error:",err)
		}
	}()
	userList := []User{}
	_ = one.All(context.TODO(), &userList)   // 当然也可以用   next
	for _, r := range userList{
		fmt.Println(r)
	}

	// delete
	// collection.DeleteOne()
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{ "name" : "bluefrog" })
	if err != nil {
		log.Fatal("DeleteMany error: ",err)
	}
	fmt.Println("delete result: ", deleteResult) // &{1}
	fmt.Printf("DeletedCount: %v\n", deleteResult.DeletedCount) // 1

	//如果我们不在使用 链接对象，那最好断开，减少资源消耗
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("disconnect error: ",err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
