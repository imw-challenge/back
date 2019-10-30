package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/imw-challenge/back/types"
)

func TestDB(t *testing.T) {
	//Create memdb and load csv
	mdb, err := InitMessageDB()
	if err != nil {
		t.Errorf("Error initializing database: %s", err)
	}

	time1, _ := time.Parse(time.RFC3339, "2017-12-14T06:20:33-08:00")
	_, offset1 := time1.Zone()
	time2, _ := time.Parse(time.RFC3339, "2018-11-24T13:16:07-08:00")
	_, offset2 := time2.Zone()
	// Insert messages
	testMessages := []*types.Message{
		&types.Message{"2C7BCEC7-CD14-D6E5-3FBF-F9551375429A", "Alex Mustermann", "fake@site.biz", "get the message?", time1.Unix(), offset1},
		&types.Message{"B4F7A417-424E-2B99-87B6-5CA0744B7BBD", "Reggie Tester", "never@existed.io", "lo and behold", time2.Unix(), offset2},
	}

	err = mdb.InsertMessages(testMessages)
	if err != nil {
		t.Errorf("Error committing one or more messages: %s", err)
	}

	message, err := mdb.FetchByID("2C7BCEC7-CD14-D6E5-3FBF-F9551375429A")
	if err != nil {
		t.Errorf("Error fetching message 2C7BCEC7-CD14-D6E5-3FBF-F9551375429A: %s", err)
	}

	fmt.Printf("Hello %s!\n", message.Name)
	if message.Name != "Alex Mustermann" {
		t.Errorf("Error fetching correct message. Expected %v, got %v", testMessages[0], message)
	}

	it, err := mdb.FetchByTime()
	for obj := it.Next(); obj != nil; obj = it.Next() {
		m := obj.(*types.Message)
		messageLocation := time.FixedZone("", m.TZ)
		zuluTime := time.Unix(m.Time, 0)
		//		localTime, _ := time.ParseInLocation(time.RFC3339, zuluTime.Format(time.RFC3339), messageLocation)
		fmt.Printf("message from %s at %s\n", m.Name, zuluTime.In(messageLocation).Format(time.RFC3339))
	}
}
