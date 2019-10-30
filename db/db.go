package db

import (
	"local/back/types"

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
						Indexer: &memdb.UUIDFieldIndex{Field: "ID"},
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

// Insert creates if messages does not exist, updates if it does exist
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

func (m *MessageDB) FetchByID(ID string) (*types.Message, error) {
	// Create read-only transaction
	txn := m.db.Txn(false)
	defer txn.Abort()

	// Lookup by message id
	raw, err := txn.First("message", "id", ID)
	if err != nil {
		return nil, err
	}
	return raw.(*types.Message), nil
}

func (m *MessageDB) FetchByTime() (ResultIter, error) {
	txn := m.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("message", "time")
	if err != nil {
		return it, err
	}
	return it, nil

}

/*


// Fetch by Name
it, err := txn.Get("message", "name", "Alex Mustermann")
if err != nil {
	panic(err)
}

fmt.Println("Match:")
for obj := it.Next(); obj != nil; obj = it.Next() {
	m := obj.(*types.Message)
	fmt.Printf("%s %s %s %s %s\n", m.ID, m.Name, m.Email, m.Text, time.fromUnix(m.UTime).tostring())
}

// Fetch messages by Time

*/
