package types

type Message struct {
	ID    string
	Name  string
	Email string
	Text  string
	Time  int64 //Unix Epoch Seconds
	TZ    int   //Seconds East of UTC
}
