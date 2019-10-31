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

//Alias Message so that we can inherit fields without inheriting methods
//   to avoid looping on Unmarshall
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
