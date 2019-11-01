package db

import (
	"sort"
	"testing"

	"github.com/imw-challenge/back/types"
)

func initEmptyDB() *MessageDB {
	mdb, err := InitMessageDB()
	if err != nil {
		//		t.Errorf("Error initializing database: %s", err)
	}
	return mdb
}

func getTestMessages() []*types.Message {
	messageStrings := []string{
		`{"id":"A5D00000-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi there","time":"2015-11-05T13:15:17-07:00"}`,
		`{"id":"B4F7A417-424E-2B99-87B6-5CA0744B7BBD","name":"Reggie Tester","email":"false@email.address","text":"lorem ipsum dolor sit amet","time":"2016-04-10T15:15:17-07:00"}`,
		`{"id":"2C7BCEC7-CD14-D6E5-3FBF-F9551375429A","name":"Alex Mustermann","email":"fake@site.biz","text":"testing","time":"2017-05-30T15:26:38-07:00"}`,
		`{"id":"B5D11111-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi again","time":"2018-09-10T00:15:00-07:00"}`,
		`{"id":"C5D22222-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"another one","time":"2019-04-10T10:10:10-07:00"}`}
	var testMessages []*types.Message
	for _, m := range messageStrings {
		msg := new(types.Message)
		msg.UnmarshalJSON([]byte(m))
		testMessages = append(testMessages, msg)
	}
	return testMessages
}

func initPopulatedDB() *MessageDB {
	mdb := initEmptyDB()
	messages := getTestMessages()
	mdb.InsertMessages(messages)
	return mdb
}

func TestInsert(t *testing.T) {
	mdb := initEmptyDB()
	testMessages := getTestMessages()

	err := mdb.InsertMessages(testMessages)
	if err != nil {
		t.Errorf("Error committing one or more messages: %s", err)
	}

	contents, err := mdb.FetchAll()
	if err != nil {
		t.Errorf("Error fetching results: %s", err)
	}

	if len(contents) != len(testMessages) {
		t.Errorf("Expected to find %d messages in db, found %d", len(testMessages), len(contents))
	}
	//get all, check count?
}

func TestFetchByID(t *testing.T) {
	mdb := initPopulatedDB()

	message, err := mdb.FetchByID("2C7BCEC7-CD14-D6E5-3FBF-F9551375429A")
	if err != nil {
		t.Errorf("Error fetching message 2C7BCEC7-CD14-D6E5-3FBF-F9551375429A: %s", err)
	}

	if message.Name != "Alex Mustermann" {
		t.Error("Error fetching correct message by ID")
	}
}

func TestFetchAntiChrono(t *testing.T) {
	mdb := initPopulatedDB()

	//Fetch from one second after earliest test message, to one second before latest
	messages, err := mdb.FetchSortedByTime(1446754518, 1554916209, false)
	if err != nil {
		t.Errorf("Error fetching messages by time: %s", err)
	}

	if len(messages) != 3 {
		t.Errorf("Fetched incorrect number of results, expected 3, got %d", len(messages))
	}
	var times []int64
	for _, m := range messages {
		//push times to front of slice
		times = append([]int64{m.Time}, times...)
	}

	if !sort.IsSorted(types.Int64Slice(times)) {
		t.Errorf("Results not sorted reverse chronologically")
	}
}
