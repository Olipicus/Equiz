package equiz

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestRegisterEvent(t *testing.T) {

	service := New("../config.json")

	mgh := service.mgh.GetSession()
	defer mgh.Close()

	collEvent := mgh.DB(service.mgh.DBName).C("event")
	collUserEvent := mgh.DB(service.mgh.DBName).C("user_event")

	user := User{ID: bson.NewObjectId(), LineID: "test01", UserName: "testuser", Pic: "test.png"}
	event := Event{ID: bson.NewObjectId(), EventTag: "#test", EventName: "test event"}

	//Prepare Event
	if err := collEvent.Insert(event); err != nil {
		t.Errorf("Insert Event Error: %v", err)
	}

	//Test Not Found Event
	if err := service.RegisterEvent(&user, &Event{EventTag: "#xxx"}); err != errorNotFoundEvent {
		t.Errorf("It should not found event error %v", err)
	}

	//Test Register Event
	if err := service.RegisterEvent(&user, &event); err != nil {
		t.Errorf("Got error on register event $%v", err)
	}

	//Test User Repeat register event
	if err := service.RegisterEvent(&user, &event); err != errorUserExist {
		t.Errorf("Should got user exist %v", err)
	}

	collEvent.Remove(bson.M{"_id": event.ID})
	collUserEvent.Remove(bson.M{"_id": user.ID})
}
