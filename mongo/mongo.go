/**
 * @Author KYIMH
 * @Description
 * @Date 2021/7/27 20:20
 **/

package mongo

import (
	"context"
	"github.com/KYIMH/CCS_Utils/share/enum/staict_const"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

//Cli -> client of mongo db
//Client: mongo client
//Database: mongo database
//Coll: mongo collection
//RetryTimes: try to ping mongodb RetryTimes times
//Timeout: the timeout of mongodb ping
//Config: include Uri, Database, Coll, Auth params
type Cli struct {
	Ctx        context.Context
	Client     *qmgo.Client
	Database   *qmgo.Database
	Coll       *qmgo.Collection
	RetryTimes int
	Timeout    int64
	Config     *qmgo.Config
}

//MgoConfig -> mongo config read from zk data
//DbNAme: name of the database you want to connect in mongodb
//Coll: name of the database you want to handle under previous database
//Uri: address of mongodb example: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
//RetryTimes: refer to the explanation above (Cli.RetryTimes)
//Timeout: refer to the explanation above (Cli.TimeOut)
//AuthDb: the name of the database to use for authentication, necessary if user is not 'admin' user
type MgoConfig struct {
	DbName     string `json:"db_name"`
	Coll       string `json:"coll"`
	Uri        string `json:"uri"`
	RetryTimes int    `json:"retry_times"`
	Timeout    int64  `json:"timeout"`
	Username   string `json:"user_name"`
	Password   string `json:"password"`
	AuthDb     string `json:"auth_db"`
}

type (
	chatMsgType staict_const.ChatMsg
	bsonM       bson.M
)

//operators of mongo client like: create mongo client, init mongo pool ...
type MogClient interface {
	initPool()
	AddClient2Pool(mongoConfig MgoConfig) error
	GetClient(dbName string) (*Cli, error)
	CreateFixedMongoCli(config MgoConfig) (*Cli, error)
	GetCtx() context.Context
	Close() error
}

// mongo data operators
type MogDal interface {
	//Create
	InsertDoc(dbName string, data interface{}) (*qmgo.InsertOneResult, error)

	//Retrieve
	GetDoc(dbName string, condition bsonM, res chatMsgType) error

	//Update
	UpdateDoc(dbName string, condition bsonM, operator bsonM) error

	//Delete
	RemoveDoc(dbName string, condition bsonM) error
}
