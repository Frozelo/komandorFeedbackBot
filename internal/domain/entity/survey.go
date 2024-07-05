package entity

type Survey struct {
	Id        int
	UserId    int
	Questions []Question
	AvgScore  float64
}

type Question struct {
	Id       int
	SurveyId int
	Text     string
	Answer   int
}
