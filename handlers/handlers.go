package handlers

import (
	db "../repositories"
	model "../models"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)


var (
	Quizz model.Quizz
 	LenQuizz int//số câu hỏi cần trả lời trong quizz
	Index int = 0//chỉ số của câu hỏi dùng cho việc chuyển câu hỏi, hiển thị
	UserLogin model.User//user đã login
)

//server push
func ServerPush(context *gin.Context){
	if pusher := context.Writer.Pusher(); pusher != nil {
		//server push
		pusher.Push("static/css/login.css", nil)
		pusher.Push("static/css/question.css", nil)
		pusher.Push("static/js/alert.js", nil)
	}
}

//render home
func GetHome(context *gin.Context){
	ServerPush(context)//push css và js
	context.HTML(http.StatusOK, "login-form.tmpl", gin.H{})
}

//render template đăng ký
func GetRegister(context *gin.Context){
	ServerPush(context)//push css và js
	context.HTML(http.StatusOK, "register.tmpl", gin.H{})
}

//render template đăng nhập thành công
func LoginSuc(context *gin.Context){
	ServerPush(context)//push css và js
	context.HTML(http.StatusOK, "loginSucceed.tmpl", gin.H{"message" : "Login succeed!", "lenquizz" : LenQuizz})
}

//render template đăng ký thành công
func RegistSuc(context *gin.Context){
	ServerPush(context)//push css và js
	context.HTML(http.StatusOK, "registSucceed.tmpl", gin.H{"message" : "Regist succeed!",})
}

//xử lý đăng ký
func PostRegister(context *gin.Context){
	ServerPush(context)//push css và js

	var userTmp model.User//biến tạm nhận kiểu user
	context.ShouldBind(&userTmp)//lấy username và password từ post form 
	
	check := false
	if userTmp.Password == "" || userTmp.Username == "" {
		check = true//nếu rỗng thì lỗi
	}
	for i := 0; i < len(userTmp.Username); i++ {
		if(userTmp.Username[i]) == ' '{
			check = true//có dấu cách trong user name thì lỗi
		}
	}

	//sinh id cho user mới
	userTmp.IdU = bson.NewObjectId()
	
	collection := db.Session.DB("assignment").C("user")//truy cập bảng employee trong db mydb
	
	var result bson.M
	//kiểm tra xem tồn tại trong db chưa
	errUname := collection.Find(bson.M{"username" : userTmp.Username}).One(&result)

	switch{
	case errUname == nil://user name đã tồn tại thì
		log.Println("a register has been refused: user exist")
		context.HTML(http.StatusOK, "registSucceed.tmpl", gin.H{"message" : "Username đã tồn tại!",})//thông báo user đã tồn tại

	case check://trong username có dấu cách hoặc rỗng thì
		log.Println("a register has been refused: user contain space charactor or empty")
		context.HTML(http.StatusOK, "registSucceed.tmpl", gin.H{"message" : "Username không được chứa khoảng trống",})//thông báo lỗi định dạng username

	default://không có lỗi thì
		userTmp.Point = 0
		err := collection.Insert(&userTmp)//chèn một record mới với thông tin user đăng ký
		if err != nil {
			log.Print("error occurred while inserting document in database :: ", err)
			log.Print("a register has been refused: cant inser new record")
			//quay lai trang dang ky
			GetRegister(context)//quay lại trang dăng ký
			return
		}else{
			log.Println("A register success")
			RegistSuc(context)//chuyển sang trang đăng nhập
			return
		}
	}
}

//xử lý đăng nhập:
func PostHome(context *gin.Context){
	
	context.ShouldBind(&UserLogin)//lấy username và password từ post form cho vào

	//truy cập bảng employee trong db mydb
	collection := db.Session.DB("assignment").C("user")
	
	var result bson.M
	//kiểm tra xem tồn tại user không
	errUser := collection.Find(bson.M{"username" : UserLogin.Username, "password" : UserLogin.Password}).One(&result)
	if errUser == nil{
		//mỗi lần đăng nhập sẽ reset điểm và index câu hỏi
		Index = 0
		UserLogin.Point = 0
		log.Print("a login Success")//Đăng nhập thành công
		
		GetQuizz(context)//cho câu hỏi vào quizz và lấy giá trị độ dài quizz
		LoginSuc(context)//render đăng nhập thành công và số câu hỏi cần phải làm
		RenderQuestion(context, Quizz[0])//render câu hỏi đầu tiên
		return
	}else{
		log.Print("a login Failed")//đăng nhập thất bại
		http.Redirect(context.Writer, context.Request, "/", http.StatusSeeOther)//quay lại trang dăng nhập
		return
	}
}

//lấy dữ liệu data cho vào quizz và lấy độ dài của quizz
func GetQuizz(context *gin.Context){
	log.Print("reading question from database")
	
	collection := db.Session.DB("assignment").C("question")
	//lấy dữ liệu từ db cho vào quizz
	err := collection.Find(bson.M{}).All(&Quizz)
	if err != nil {
		log.Print("error occurred while reading documents from database :: ", err)
		return
	}	
	LenQuizz = len(Quizz)//lấy giá trị độ dài của quizz
}

//render câu hỏi
func RenderQuestion(context *gin.Context, question model.QuestionType){
	ServerPush(context)
	context.HTML(http.StatusOK, "question.tmpl", bson.M{
		"Id" : question.Id,
		"Type" : question.Type,
		"Question" : question.Question,
		"Answer1" : question.Answer1.Answer,
		"Answer2" : question.Answer2.Answer,
		"Answer3" :question.Answer3.Answer,
		"Answer4" : question.Answer4.Answer,
		"Correct1" : question.Answer1.Correct,
		"Correct2" : question.Answer2.Correct,
		"Correct3" : question.Answer3.Correct,
		"Correct4" : question.Answer4.Correct,
		//2 giá trị ở dưới để xử lý logic trong template
		"lenquizz": LenQuizz - 1,
		"index" : Index,
	})
}

//submit câu trả lời và update điểm
func Submit(context*gin.Context){
	var canswer  model.Answer//lưu dữ liệu trả lời của câu hỏi checkbox
	var ranswer model.Answer2//lưu dữ liệu trả lời của câu hỏi radio
	

	check := false//check xem có đáp án nào sai không, mặc định là không
	
	if Quizz[Index].Type == "checkbox" {
		context.ShouldBind(&canswer)//doc vao dap an checbox
		if Quizz[Index].Answer1.Correct && canswer.AnswerC1 == ""{
			check = true
		}
		if !Quizz[Index].Answer1.Correct && canswer.AnswerC1 == "true"{
			check = true
		}
		//kiem tra dap an checkbox2
		if Quizz[Index].Answer2.Correct && canswer.AnswerC2 == ""{
			check = true
		}
		if !Quizz[Index].Answer2.Correct && canswer.AnswerC2 == "true"{
			check = true
		}
		//kiem tra dap an checkbox3
		if Quizz[Index].Answer3.Correct && canswer.AnswerC3 == ""{
			check = true
		}
		if !Quizz[Index].Answer3.Correct && canswer.AnswerC3 == "true"{
			check = true
		}
		//kiem tra dap an checkbox4
		if Quizz[Index].Answer4.Correct && canswer.AnswerC4 == ""{
			check = true
		}
		if !Quizz[Index].Answer4.Correct && canswer.AnswerC4 == "true"{
			check = true
		}
		
	}else{
		context.ShouldBind(&ranswer)//doc vao dap an radio
		if ranswer.Answer == "false" || ranswer.Answer == ""{
			check = true
		}
	}

	//kiểm tra đáp án và tăng điểm
	if !check{//nếu không có đáp án nào sai
		UserLogin.Point ++//tăng điểm
	}

	log.Print(UserLogin.Point)
	log.Print("Answer checkbox: ", canswer)
	log.Print("Answer radio: ", ranswer)

	//log.Print("radio answer: ", ranswer)
	Index++
	if Index < LenQuizz {
		RenderQuestion(context, Quizz[Index])
	}
	
}

//type nhận vào button
type Button struct{
	Event string `form:"switchQ"`
}
//xử lý các button chức năng trong quizz
func SwitchQuestion(context *gin.Context)  {
	var event Button
	context.ShouldBind(&event)
	switch event.Event{
	case "previous"://neu button previous duoc nhan
		Index--
		RenderQuestion(context, Quizz[Index])//render câu hỏi Index

	case "next"://neu button next duoc nhan
		Index++
		RenderQuestion(context, Quizz[Index])//render câu hỏi Index
	
	case "sub"://sumit dap an
		Submit(context)
		
	case "end"://render ra diem quizz
		Submit(context)
		context.HTML(http.StatusOK, "point.tmpl", gin.H{"point" : UserLogin.Point, "lenquizz" : LenQuizz})
	}
}


