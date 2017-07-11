package mongodb

import (
	"errors"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Helper Struct of MongoHelper
type Helper struct {
	session *mgo.Session
	DBName  string
}

//Init : Initial DB
func (h *Helper) Init(mongoAddress string, databaseName string) {
	session, err := mgo.Dial(mongoAddress)
	h.DBName = databaseName
	if err != nil {
		log.Fatal(err)
	}
	h.session = session
	//log.Println("Connect MongoDB OK!")

	// Optional. Switch the session to a monotonic behavior.
	h.session.SetMode(mgo.Monotonic, true)
}

//Init : Initial DB
func (h *Helper) InitWithAuth(mongoAddress string, databaseName string, username string, password string) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoAddress},
		Database: databaseName,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	h.DBName = databaseName
	if err != nil {
		log.Fatal(err)
	}
	h.session = session
	//log.Println("Connect MongoDB OK!")

	// Optional. Switch the session to a monotonic behavior.
	h.session.SetMode(mgo.Monotonic, true)
}

//GetSession : Get DB Session
func (h *Helper) GetSession() *mgo.Session {
	return h.session.Clone()
}

//GetOneData : Get Single Document
func (h *Helper) GetOneData(collectionName string, id interface{}) (interface{}, error) {
	session := h.session.Clone()
	defer session.Close()
	c := h.session.DB(h.DBName).C(collectionName)
	var obj interface{}
	var err error

	switch id.(type) {
	case string:
		err = c.FindId(id.(string)).One(&obj)
	case bson.ObjectId:
		err = c.FindId(id.(bson.ObjectId)).One(&obj)
	}

	return obj, err
}

//GetOneDataToObj : Get Single Document
func (h *Helper) GetOneDataToObj(collectionName string, id interface{}, obj interface{}) error {
	session := h.session.Clone()
	defer session.Close()
	c := session.DB(h.DBName).C(collectionName)
	var err error

	switch id.(type) {
	case string:
		err = c.FindId(id.(string)).One(obj)
	case bson.ObjectId:
		err = c.FindId(id.(bson.ObjectId)).One(obj)
	}
	return err
}

//RemoveByID : Remove Data By ID
func (h *Helper) RemoveByID(collectionName string, id interface{}) error {
	session := h.session.Clone()
	defer session.Close()
	c := session.DB(h.DBName).C(collectionName)
	switch id.(type) {
	case string:
		return c.RemoveId(id.(string))
	case bson.ObjectId:
		return c.RemoveId(id.(bson.ObjectId))
	}
	return errors.New("Type of ID Mismatch")
}

//InsertData : Insert Document to Collection
func (h *Helper) InsertData(collectionName string, obj interface{}) error {
	session := h.session.Clone()
	defer session.Close()
	c := session.DB(h.DBName).C(collectionName)
	return c.Insert(obj)
}

//UpdateData : Update Document
func (h *Helper) UpdateData(collectionName string, id interface{}, obj interface{}) error {
	session := h.session.Clone()
	defer session.Close()
	c := session.DB(h.DBName).C(collectionName)
	update := bson.M{"$set": obj}
	switch id.(type) {
	case string:
		return c.UpdateId(id.(string), update)
	case bson.ObjectId:
		return c.UpdateId(id.(bson.ObjectId), update)
	}
	return errors.New("Type of ID Mismatch")
}
