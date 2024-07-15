package entities

type Request struct {
	RemoteAddr string      `json:"-" bson:"remote_addr"`
	Request    PollRequest `json:"request" bson:"request" validate:"required,required"`
}
