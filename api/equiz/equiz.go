package equiz

import (
	"encoding/json"
	"errors"

	"code.olipicus.com/equiz/config"
	"code.olipicus.com/equiz/db/mongodb"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	LineID   string        `json:"line_id" bson:"line_id"`
	UserName string        `json:"user_name" bson:"user_name"`
	EventTag string        `json:"event_tag" bson:"event_tag"`
	Pic      string        `json:"pic" bson:"pic"`
}

type Event struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	EventTag  string        `json:"event_tag" bson:"event_tag"`
	EventName string        `json:"event_name" bson:"event_name"`
}

type EquizService struct {
	mgh        *mongodb.Helper
	configPath string
	config     *config.Table
}

var ErrorUserExist error = errors.New("User is exists")
var ErrorNotFoundEvent error = errors.New("Not Found Event")

func New(configPath string) *EquizService {

	if configPath == "" {
		configPath = "config.json"
	}

	service := EquizService{}
	service.configPath = configPath
	service.config = config.LoadConfig(configPath)

	b, _ := json.Marshal(service.config.Get("mongodb"))
	var mongoConfig mongodb.MongoConfig
	json.Unmarshal(b, &mongoConfig)

	service.mgh = &mongodb.Helper{}
	service.mgh.Init(mongoConfig.Address, mongoConfig.DB)
	return &service
}

func (service *EquizService) RegisterEvent(user *User, event *Event) error {
	mghSession := service.mgh.GetSession()
	defer mghSession.Close()

	collEvent := mghSession.DB(service.mgh.DBName).C("event")
	countEvent, err := collEvent.Find(bson.M{"event_tag": event.EventTag}).Count()

	if countEvent == 0 {
		return ErrorNotFoundEvent
	}

	if err != nil {
		return err
	}

	user.EventTag = event.EventTag
	collUserEvent := mghSession.DB(service.mgh.DBName).C("user_event")
	countUser, err := collUserEvent.Find(bson.M{"event_tag": event.EventTag, "line_id": user.LineID}).Count()

	switch {
	case countUser == 0:
		collUserEvent.Insert(user)
		return nil
	case countUser >= 1:
		return ErrorUserExist
	}

	return nil

}
