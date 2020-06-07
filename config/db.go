package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DBClient used in other package
	DBClient *mongo.Client
	// DB other package
	DB *mongo.Database
	dc dbConfig
)

type dbConfig struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
}

func init() {
	getDbConfigInfo()
}

func getDbConfigInfo() {
	var err error
	configDir := os.Getenv("ConfigDir")
	dbConfigPath := path.Join(configDir, "dbconfig.toml")
	_, err = os.Stat(dbConfigPath)
	if os.IsNotExist(err) {
		f, _ := os.OpenFile(dbConfigPath, os.O_CREATE|os.O_RDWR, 0644)
		defaultContent :=
			`
DbHost="127.0.0.1"
DbPort="27017"
DbName="fiberapp"
DbUser="root"
DbPassword="password"
`
		f.WriteString(defaultContent)
		f.Close()
		toml.Decode(defaultContent, &dc)
		return
	}
	toml.DecodeFile(dbConfigPath, &dc)
}

// DBConnect a mongodb
func DBConnect() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://" + dc.DbHost + ":" + dc.DbPort)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DBClient, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = DBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DB = DBClient.Database(dc.DbName)
	fmt.Println("Connected to MongoDB!")
}

// DBDisconnect mongodb
func DBDisconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	DBClient.Disconnect(ctx)
}
