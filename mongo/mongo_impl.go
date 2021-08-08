/**
 * @Author KYIMH
 * @Description
 * @Date 2021/7/27 20:20
 **/

package mongo

import (
	"context"
	"errors"
	"github.com/KYIMH/CCS_Utils/share/enum/staict_const"
	"github.com/qiniu/qmgo"
)

type MogPoolType map[string]*Cli

//MogClientImpl -> mongo client implement
//Config: list of MogConfig
//Pool: map of Cli(mongo client) example: {'dbname': mongo client of dbname}
type MogClientImpl struct {
	Config  []MgoConfig
	Pool    MogPoolType
	Context context.Context
}

//create new mongodb client
func NewMongoClient() *MogClientImpl {

	cli := new(MogClientImpl)

	cli.initPool()

	cli.Context = context.Background()

	return cli
}

//init mongo pool, and pool will be empty
func (m *MogClientImpl) initPool() {

	m.Pool = make(MogPoolType, 0)
}

//add new mongo client to connection pool
func (m *MogClientImpl) AddClient2Pool(mongoConfig MgoConfig) error {

	m.Config = append(m.Config, mongoConfig)

	for _, config := range m.Config {

		newCli, err := m.CreateFixedMongoCli(config)

		if nil != err {
			return err
		}
		m.Pool[config.DbName] = newCli
	}

	return nil
}

//get mongo client by dbname
func (m *MogClientImpl) GetClient(dbName string) (*Cli, error) {

	cli, ok := m.Pool[dbName]

	if !ok {
		return nil, errors.New("no connection " + dbName + " in Manager")
	}
	return cli, nil
}

// return context for mongodb
func (m *MogClientImpl) GetCtx() context.Context {

	return m.Context
}

//use this function to get a connection points to a fixed database and collection
func (m *MogClientImpl) CreateFixedMongoCli(config MgoConfig) (*Cli, error) {

	connConfig := &qmgo.Config{
		Uri:      config.Uri,
		Database: config.DbName,
		Coll:     config.Coll,
		Auth: &qmgo.Credential{
			Username:   config.Username,
			Password:   config.Password,
			AuthSource: config.AuthDb,
		},
	}

	cli, err := qmgo.Open(m.Context, connConfig)

	if nil != err {
		return nil, err
	}

	//create new mongo client
	newCli := &Cli{
		Client:     cli.Client,
		Database:   cli.Database,
		Coll:       cli.Collection,
		RetryTimes: config.RetryTimes,
		Timeout:    config.Timeout,
		Config:     connConfig,
	}

	return newCli, nil
}

//close all mongo client
func (m *MogClientImpl) Close() error {

	for _, cli := range m.Pool {
		err := cli.Client.Close(m.Context)
		if nil != err {
			return err
		}
	}

	return nil
}

/*================
Mongo dal operator
==================*/

//insert one document
func (m *MogClientImpl) InsertDoc(dbName string, data interface{}) (*qmgo.InsertOneResult, error) {

	if "" == dbName {
		dbName = staict_const.Chat
	}

	cli, err := m.GetClient(dbName)
	if nil != err {
		return nil, err
	}

	result, err := cli.Coll.InsertOne(m.GetCtx(), &data)
	if nil != err {
		return nil, err
	}

	return result, nil
}

//get one document
func (m *MogClientImpl) GetDoc(dbName string, condition bsonM, res chatMsgType) error {

	if "" == dbName {
		dbName = staict_const.Chat
	}

	cli, err := m.GetClient(dbName)
	if nil != err {
		return err
	}

	err = cli.Coll.Find(m.GetCtx(), condition).One(&res)

	if nil != err {
		return err
	}

	return nil
}

//update one document
func (m *MogClientImpl) UpdateDoc(dbName string, condition bsonM, operator bsonM) error {

	if "" == dbName {
		dbName = staict_const.Chat
	}

	cli, err := m.GetClient(dbName)
	if nil != err {
		return err
	}

	//example: err = cli.Coll.UpdateOne(dal.Client.GetCtx(), bson.M{"name": "d4"}, bson.M{"$set": bson.M{"age": 7}})
	err = cli.Coll.UpdateOne(m.GetCtx(), condition, operator)

	if nil != err {
		return err
	}

	return nil
}

//remove one doc
func (m *MogClientImpl) RemoveDoc(dbName string, condition bsonM) error {

	if "" == dbName {
		dbName = staict_const.Chat
	}

	cli, err := m.GetClient(dbName)
	if nil != err {
		return err
	}

	err = cli.Coll.Remove(m.GetCtx(), condition)
	if nil != err {
		return err
	}

	return nil
}
