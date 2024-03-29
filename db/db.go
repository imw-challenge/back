package db

import (
	"encoding/csv"
	"errors"
	"os"
	"sort"
	"time"

	"github.com/imw-challenge/back/types"

	"github.com/hashicorp/go-memdb"
)

type MessageDB struct {
	db *memdb.MemDB
}

type ResultIter memdb.ResultIterator

func InitMessageDB() (*MessageDB, error) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"message": &memdb.TableSchema{
				Name: "message",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"time": &memdb.IndexSchema{
						Name:    "time",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Time"},
					},
				},
			},
		},
	}

	// Create a new data base
	mdb, err := memdb.NewMemDB(schema)
	if err != nil {
		return &MessageDB{}, err
	}
	return &MessageDB{mdb}, nil
}

func (m *MessageDB) LoadFromCSV(filename string, batchSize int) error {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	//Remove first (label) line
	lines = lines[1:]

	var batch []*types.Message
	for _, line := range lines {
		messageTime, _ := time.Parse(time.RFC3339, line[4])
		_, timeOffset := messageTime.Zone()
		message := &types.Message{ID: line[0], Name: line[1], Email: line[2], Text: line[3], Time: messageTime.Unix(), TZ: timeOffset}
		batch = append(batch, message)
		if len(batch)%batchSize == 0 {
			err := m.InsertMessages(batch)
			if err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		err = m.InsertMessages(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

// Insert creates if message does not exist, updates if it does exist
func (m *MessageDB) InsertMessages(messages []*types.Message) error {
	// Create a write transaction
	txn := m.db.Txn(true)

	for _, m := range messages {
		if err := txn.Insert("message", m); err != nil {
			return err
		}
	}

	// Commit the transaction
	txn.Commit()
	return nil
}

// InsertMessage creates a message if it does not exist, or updates if it does exist
func (m *MessageDB) InsertMessage(message *types.Message) error {
	// Create a write transaction
	txn := m.db.Txn(true)

	if err := txn.Insert("message", message); err != nil {
		return err
	}

	// Commit the transaction
	txn.Commit()
	return nil
}

// FetchAll returns a slice with all the messages in the database, in non-deterministic order
func (m *MessageDB) FetchAll() ([]*types.Message, error) {
	txn := m.db.Txn(false)
	it, err := txn.Get("message", "id")
	if err != nil {
		return []*types.Message{}, err
	}

	var messages []*types.Message
	for obj := it.Next(); obj != nil; obj = it.Next() {
		m := obj.(*types.Message)
		messages = append(messages, m)
	}
	return messages, nil
}

// FetchByID fetches a single message by ID
func (m *MessageDB) FetchByID(ID string) (*types.Message, error) {
	// Create read-only transaction
	txn := m.db.Txn(false)
	defer txn.Abort()

	// Lookup by message id
	raw, err := txn.First("message", "id", ID)
	if err != nil {
		return nil, err
	}
	if raw != nil { //message found
		return raw.(*types.Message), nil
	} else {
		return &types.Message{}, errors.New("Message not found")
	}
}

// fetchByMinTime returns a result iterator for messages with a timestamp
// at or after the specified start in unix seconds
func (m *MessageDB) fetchByMinTime(start int64) (ResultIter, error) {
	txn := m.db.Txn(false)
	defer txn.Abort()

	it, err := txn.LowerBound("message", "time", start)
	if err != nil {
		return it, err
	}
	return it, nil

}

// FetchSortedByTime returns a slice of messages with timestamps between
// start and end (inclusive), specified in unix seconds - the ordering of this array
// is chronological if ascending is true, or reverse chronological if false
func (m *MessageDB) FetchSortedByTime(start int64, end int64, ascending bool) ([]*types.Message, error) {
	//fetch by time
	var messages []*types.Message
	it, err := m.fetchByMinTime(start)
	if err != nil {
		return []*types.Message{}, err
	}

	//insert into slice
	for obj := it.Next(); obj != nil; obj = it.Next() {
		m := obj.(*types.Message)
		if m.Time <= end {
			messages = append(messages, m)
		}
	}

	//sort slice
	if ascending {
		sort.Slice(messages, func(i, j int) bool { return messages[i].Time < messages[j].Time })
	} else {
		sort.Slice(messages, func(i, j int) bool { return messages[i].Time > messages[j].Time })
	}
	return messages, nil
}
