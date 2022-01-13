package repositories

import(
	mgo "gopkg.in/mgo.v2"
	"log"
)

const (
	CONN_HOST    = "localhost"
	MONGO_DB_URL = "mongodb://127.0.0.1:27017"
)
var ConnectionError error
var Session *mgo.Session//tạo 1 session vào db

//kết nối tới database
func init() {
	Session, ConnectionError = mgo.Dial(MONGO_DB_URL)
	if ConnectionError != nil {
		log.Fatal("error connecting to database :: ", ConnectionError)
	}
	Session.SetMode(mgo.Monotonic, true)
}