package entity

type Survey struct {
	Questions []Question
}

type Question struct {
	Id       int
	Text     string
	Answered bool
	Answer   string
}
