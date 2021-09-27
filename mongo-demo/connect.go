package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

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


	//如果我们不在使用 链接对象，那最好断开，减少资源消耗
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("disconnect error: ",err)
	}
	fmt.Println("Connection to MongoDB closed.")
}