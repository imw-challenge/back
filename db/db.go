package db

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
)

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
db, err := memdb.NewMemDB(schema)
if err != nil {
	panic(err)
}

// Create a write transaction
txn := db.Txn(true)

// Insert messages
messages := []*types.Message{
	&types.Message{"2C7BCEC7-CD14-D6E5-3FBF-F9551375429A", "Alex Mustermann", "fake@site.biz", "get the message?", "2017-12-14T06:20:33-08:00"},
	&types.Message{"B4F7A417-424E-2B99-87B6-5CA0744B7BBD", "Reggie Tester", "never@existed.io", "lo and behold", "2018-11-24T13:16:07-08:00"},
}
for _, m := range messages {
	if err := txn.Insert("message", m); err != nil {
		panic(err)
	}
}

// Commit the transaction
txn.Commit()

// Create read-only transaction
txn = db.Txn(false)
defer txn.Abort()

// Lookup by message id
raw, err := txn.First("message", "id", "B4F7A417-424E-2B99-87B6-5CA0744B7BBD")
if err != nil {
	panic(err)
}

// Say hi!
fmt.Printf("Hello %s!\n", raw.(*Message).Name)

// Fetch by ID
it, err := txn.Get("message", "id", "B4F7A417-424E-2B99-87B6-5CA0744B7BBD")
if err != nil {
	panic(err)
}

fmt.Println("Match:")
for obj := it.Next(); obj != nil; obj = it.Next() {
	m := obj.(*types.Message)
	fmt.Printf("%s %s %s %s %s\n", m.ID, m.Name, m.Email, m.Text, time.fromUnix(m.UTime).tostring())
}

// Fetch messages by Time
it, err = txn.Get("message", "time")
if err != nil {
	panic(err)
}

for obj := it.Next(); obj != nil; obj = it.Next() {
	m := obj.(*types.Message)
	fmt.Printf("message from %s at %s\n", m.Name, m.time.stringyfy)
}