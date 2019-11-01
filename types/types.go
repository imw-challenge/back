package types

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Text  string `json:"text"`
	Time  int64  //Unix Epoch Seconds
	TZ    int    //Seconds East of UTC
}

// Custom marshaller for Message, converts unix seconds + offset to RFC3339 format
func (m *Message) MarshalJSON() ([]byte, error) {
	messageLocation := time.FixedZone("", m.TZ)
	utc := time.Unix(m.Time, 0)
	messageTime := utc.In(messageLocation).Format(time.RFC3339)
	return json.Marshal(&struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Text  string `json:"text"`
		Time  string `json:"time"`
	}{
		ID:    m.ID,
		Name:  m.Name,
		Email: m.Email,
		Text:  m.Text,
		Time:  messageTime,
	})
}

// Custom unmarshaller for Message converts RFC3339 format to unix seconds + offset
// Aliases Message so that we can inherit fields without inheriting methods
//   to avoid looping on UnmarshalJSON
func (m *Message) UnmarshalJSON(data []byte) error {
	type MessageAlias Message
	aux := &struct {
		Time string `json:"time"`
		*MessageAlias
	}{
		MessageAlias: (*MessageAlias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	messageTime, _ := time.Parse(time.RFC3339, aux.Time)
	_, timeOffset := messageTime.Zone()
	m.Time = messageTime.Unix()
	m.TZ = timeOffset
	return nil
}

//Int64Slice attaches sort interface methods to []int64
//Allows for sort check at end of TestFetchAntiChrono
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
