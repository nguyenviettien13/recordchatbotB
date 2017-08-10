package main

import (
	"database/sql"
	"fmt"
	"github.com/michlabs/fbbot"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const (
	PAGEACCESSTOKEN = ""
	VERIFYTOKEN     = ""
	PORT = 2102

	DB_NAME="recordchatbotB"
	DB_USER=""
	DB_PASS="1"
	MaxSample=180
	Maxinning=5
)

var db *sql.DB
type Record struct {}

func ( record Record ) HandleMessage( bot *fbbot.Bot , msg *fbbot.Message ){
	if IsNewUser(db,msg.Sender.ID) {
		greeting := "Xin chào! "+msg.Sender.FirstName()+" Chúng tôi đang thực hiện một dự án thu thập dữ liệu ghi âm giọng nói và rất vui khi nhận được sự hợp tác của bạn"
		m := fbbot.NewTextMessage(greeting)
		bot.Send(msg.Sender,m)

		tutorialmesseger :="Hướng dẫn: Bây h tôi sẽ gửi cho bạn một đoạn text, bạn hãy đọc và ghi âm rồi gửi chúng lại cho t ôi"
		m1 := fbbot.NewTextMessage(tutorialmesseger)
		bot.Send(msg.Sender,m1)

		start:="Oki! Bây h chúng ta sẽ bắt đầu"
		m2 := fbbot.NewTextMessage(start)
		bot.Send(msg.Sender,m2)

		_,err := db.Query("INSERT INTO UserState(FbId,LastSample,Inning,Gender) VALUES(?,?,?,?)",msg.Sender.ID,0,1,msg.Sender.Gender())
		if err!=nil{
			log.Println("error when insertNewUser")
		}
		btnms :=fbbot.NewButtonMessage()
		btnms.AddPostbackButton("Xem hướng dẫn xử dụng","xemhuongdan")
		btnms.AddPostbackButton("Bắt đầu thu âm","batdauthuam")
		btnms.Text="Bạn hãy chọn một trong hai chế độ"
		bot.Send(msg.Sender,btnms)

	} else if IsAudioMessage(msg) {
		st := GetCurrentState(db,msg.Sender.ID)
		id := st + 1
		inning := GetCurrentInning(db, msg.Sender.ID)
		if !isExist(db,msg.Sender.ID,id,inning){
			_, err := db.Query("INSERT INTO Outputs(FbId, SampleId, State, Inning, UrlRecord,Gender) Value(?,?,?,?,?,?)",msg.Sender.ID,id,false,inning,msg.Audios[0].URL,msg.Sender.Gender())
			if err != nil {
				log.Println("error when Insert to outputs")
			}
		} else {
			state := GetCurrentState(db,msg.Sender.ID)
			sampleid := state+1
			stmtInsAudio , err := db.Prepare("UPDATE Outputs SET UrlRecord=? WHERE FbId= ? AND SampleId=?")
			if err != nil {
				log.Println("error when create stminsertAudio")
			}
			stmtInsAudio.Query(msg.Audios[0].URL, msg.Sender.ID,sampleid)
		}
		btnms :=fbbot.NewButtonMessage()
		btnms.AddPostbackButton("Ghi âm lại","ghiamlai")
		btnms.AddPostbackButton("Ghi âm câu tiếp theo","cautieptheo")
		btnms.Text="Bạn muốn ghi âm lại hay ghi âm câu tiếp theo"
		bot.Send(msg.Sender,btnms)
	} else {
		state := GetCurrentState(db,msg.Sender.ID)
		smlid := state+1
		inning := GetCurrentInning(db,msg.Sender.ID)
		if inning<= Maxinning {
			if smlid <= MaxSample{
				sample := GetSample(db,smlid)
				m := fbbot.NewTextMessage(sample)
				bot.Send(msg.Sender,m)
			}
		}else {
			sample:= "Bạn đã hoàn thành quá trình ghi âm"
			m:= fbbot.NewTextMessage(sample)
			bot.Send(msg.Sender,m)
		}
	}
}

func (r Record) HandlePostback(bot *fbbot.Bot, pbk *fbbot.Postback)  {
	switch pbk.Payload {
	case "xemhuongdan":
		url := "https://www.facebook.com/tiennv000"
		m := fbbot.NewTextMessage(url)
		bot.Send(pbk.Sender,m)

	case "batdauthuam":
		provincechooser := fbbot.NewButtonMessage()

		provincechooser.AddPostbackButton ("Miền Bắc","Miền Bắc")
		provincechooser.AddPostbackButton ("Miền Trung","Miền Trung")
		provincechooser.AddPostbackButton ("Miền Nam","Miền Nam")
		provincechooser.Text = "Giọng của bạn thuộc vùng miền nào???"
		bot.Send(pbk.Sender,provincechooser)

	case "cautieptheo":
		state := GetCurrentState(db,pbk.Sender.ID)
		smlid :=state+1
		ig    := GetCurrentInning(db,pbk.Sender.ID)
		if ig<= Maxinning {
			if smlid <MaxSample {
				_, err := db.Query("UPDATE Outputs SET State = ? WHERE FbId=? AND SampleId=? AND Inning=?",true,pbk.Sender.ID,smlid,ig)
				if err != nil {
					log.Println("error when update state of outputs")
				}

				_, err1 := db.Query("UPDATE UserState SET LastSample=? WHERE FbId=? ",smlid,pbk.Sender.ID)
				if err1 != nil {
					log.Println("error when update state of user ")
				}

				sample := GetSample(db,smlid+1)
				m := fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)

			}else if smlid == MaxSample {
				_, err := db.Query("UPDATE Outputs SET State=? WHERE FbId=? AND SampleId=? AND Inning=?",true,pbk.Sender.ID,smlid,ig)
				if err != nil {
					log.Println("error when update state of outputs")
				}
				_, err1 := db.Query("UPDATE UserState SET LastSample=?, Inning =? where FbId=? ",0,ig+1, pbk.Sender.ID)
				if err1 != nil {
					log.Println("error when update state of user ")
				}
				ann1 := "Ban đã hoàn thành lượt:"
				ann2 := strconv.Itoa(ig)
				announce := ann1+ann2
				m := fbbot.NewTextMessage(announce)
				bot.Send(pbk.Sender,m)
				ig:=GetCurrentInning(db,pbk.Sender.ID)
				if ig>Maxinning{
					sample := "Ban đã hoàn thành quá trình thu âm, xin cảm ơn bạn"
					m := fbbot.NewTextMessage(sample)
					bot.Send(pbk.Sender,m)
				} else {
					sample := GetSample(db,1)
					m1 := fbbot.NewTextMessage(sample)
					bot.Send(pbk.Sender,m1)
				}


			}
		}else {
			sample := "Ban đã hoàn thành quá trình thu âm, xin cảm ơn bạn"
			m := fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}

	case "ghiamlai":
		state := GetCurrentState(db,pbk.Sender.ID)
		sampleid := state+1
		ig    := GetCurrentInning(db,pbk.Sender.ID)
		if ig<=Maxinning {
			sample :=GetSample(db,sampleid)
			m:= fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}else {
			sample := "Bạn đã hoàn thành xong bài ghi âm"
			m:= fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}

	case "Miền Bắc":
		_ , err :=db.Query("UPDATE UserState SET Province =? WHERE FbId= ? ","Miền Bắc",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updateprovince user")
		}
		id := GetCurrentState(db,pbk.Sender.ID)
		sml:= id+1
		ig := GetCurrentInning(db,pbk.Sender.ID)
		if ig<= Maxinning {
			if sml <= MaxSample {
				sample := GetSample(db,sml)
				m := fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)
			}
		}else {
			sample := "Ban đã hoàn thành quá trình thu âm, xin cảm ơn bạn"
			m := fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}
	case "Miền Trung":
		_ , err :=db.Query("UPDATE UserState SET Province =? WHERE FbId= ? ","Miền Trung",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updateprovince user")
		}
		id := GetCurrentState(db,pbk.Sender.ID)
		sml:= id+1
		ig := GetCurrentInning(db,pbk.Sender.ID)
		if ig<= Maxinning {
			if sml <= MaxSample {
				sample := GetSample(db,sml)
				m := fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)
			}
		}else {
			sample := "Ban đã hoàn thành quá trình thu âm, xin cảm ơn bạn"
			m := fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}
	case "Miền Nam":
		_ , err :=db.Query("UPDATE UserState SET Province =? WHERE FbId= ? ","Miền Nam",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updateprovince user")
		}
		id := GetCurrentState(db,pbk.Sender.ID)
		sml:= id+1
		ig := GetCurrentInning(db,pbk.Sender.ID)
		if ig<= Maxinning {
			if sml <= MaxSample {
				sample := GetSample(db,sml)
				m := fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)
			}
		}else {
			sample := "Ban đã hoàn thành quá trình thu âm, xin cảm ơn bạn"
			m := fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}

	default:
		log.Println("no case in switch")
	}
}
func main() {
	//processing database
	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME )//"user:password@/dbname"
	fmt.Println("Opening connection")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	fmt.Println("checked opening connnection")
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	fmt.Println("Ping database")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("checked ping database")
	fmt.Println("added database channel")


	var r Record
	bot := fbbot.New(PORT,VERIFYTOKEN,PAGEACCESSTOKEN)
	bot.AddMessageHandler(r)
	bot.AddPostbackHandler(r)

	bot.Run()
}
func GetCurrentState(db *sql.DB,FbId string)  int {
	var lastsample int
	row , err :=db.Query("SELECT LastSample FROM UserState WHERE FbId=?",FbId)
	if err != nil {
		log.Println("errors query in function GetCurrentState")
	}else {
		for row.Next() {
			err := row.Scan(&lastsample)
			if err != nil {
				log.Println("errors when Scan in GetCurrentState func")
			}
		}
	}
	return lastsample
}
func GetCurrentInning(db *sql.DB,FbId string)  int {
	var inning int
	row , err :=db.Query("SELECT Inning FROM UserState WHERE FbId=?",FbId)
	if err != nil {
		log.Println("errors query in function GetCurrentInning")
	}else {
		for row.Next() {
			err := row.Scan(&inning)
			if err != nil {
				log.Println("errors when Scan in GetCurrentInning func")
			}
		}
	}
	return inning
}
func IsNewUser(db *sql.DB,FbId string) bool {
	row , err := db.Query("SELECT * FROM UserState Where FbId=?",FbId)
	if err != nil {
		log.Println("error when execute query in IsNewUser func")
	}
	return !row.Next()
}

func GetSample(db *sql.DB,Id int) string {
	row , err := db.Query("SELECT * FROM InputText WHERE Id=?",Id)
	var id int
	var Sample string
	if err != nil {
		log.Println("error when getsample")
	}
	if row.Next() {
		row.Scan(&id,&Sample)
		return Sample
	} else {
		return ""
	}
}

func IsAudioMessage(msg *fbbot.Message) bool  {
	if len(msg.Audios)==0 {

		return false
	}
	log.Println("received audio")
	return true

}

func isExist(db *sql.DB,fbid string, sampleid int, inning int) bool {
	row, err := db.Query("SELECT * FROM Outputs Where FbId = ? AND SampleId = ? AND Inning = ? ",fbid,sampleid,inning)
	if err != nil {
		log.Println("error in IsExist function ")
	}
	return row.Next()
}

