package entity

type Survey struct {
	Id        int
	UserId    int
	Questions []Question
	AvgScore  float64
}

type Category struct {
	Id   int
	Name string
}

type Question struct {
	Id         int
	CategoryId int
	Text       string
}

type Answer struct {
	Id         int
	SurveyID   int
	QuestionID int
	Answer     int
}
