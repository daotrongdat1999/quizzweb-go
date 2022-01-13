package models
import(
	"gopkg.in/mgo.v2/bson"
)
//kiểu user
type User struct{
	IdU bson.ObjectId    `json:"_id" bson:"_id"`
	Username string `form:"username" json:"username" bson:"username"`
	Password string `form:"password" json:"password" bson:"password"`
	Point int `json:"point" bson:"point"`//điểm quizz
}
//kieu dap an
type AnswerType struct{
	Answer string `json:"aswer" bson:"aswer"`//đáp án
	Correct bool `json:"correct" bson:"correct"`//đáp án đúng hay sai
}

type Quizz []QuestionType//mảng quizz chứa nhiều câu hỏi
//kieu cau hoi
type QuestionType struct{
	IdQ bson.ObjectId `json:"_id" bson:"_id"`
	Id int `json:"IdAs" bson:"IdAs"`
	Type string `json:"type" bson:"type"`//kiểu 1 đáp án hay nhiều đáp án
	Question string `json:"question" bson:"question"`
	Answer1 AnswerType `json:"answer1" bson:"answer1"`
	Answer2 AnswerType ` json:"answer2" bson:"answer2"`
	Answer3 AnswerType `json:"answer3" bson:"answer3"`
	Answer4 AnswerType `json:"answer4" bson:"answer4"`
}

//kiểu câu trả lời cho radio
type Answer2 struct{
	Answer string `form:"answer"`
}
//kiểu câu trả lời checkbox
type Answer struct{
	AnswerC1 string `form:"answer1"`
	AnswerC2 string `form:"answer2"`
	AnswerC3 string `form:"answer3"`
	AnswerC4 string `form:"answer4"`
}