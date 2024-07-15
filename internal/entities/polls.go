package entities

type CreatePollRequest struct {
	PollID int  `json:"poll_id" bson:"poll_id" validate:"required"`
	Poll   Poll `json:"poll" bson:"poll" validate:"required,required"`
}

type Poll struct {
	Name     string    `json:"name" validate:"required"`
	Variants []Variant `json:"variants" bson:"variants" validate:"required,dive,required"`
}

type Variant struct {
	VariantID int    `json:"variant_id" bson:"variant_id" validate:"required"`
	Votes     int    `bson:"votes"`
	Text      string `json:"text" bson:"text" validate:"required"`
}

type PollRequest struct {
	PollID    int `json:"poll_id" bson:"poll_id" validate:"required"`
	VariantID int `json:"variant_id" validate:"required"`
}

type GetResultRequest struct {
	PollID int `json:"poll_id" bson:"poll_id" validate:"required"`
}

type DeleteRequest struct {
	PollID int `json:"poll_id" bson:"poll_id" validate:"required"`
}
