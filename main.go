package main

import (
	"database/sql"
	"fmt"
	"github.com/michlabs/fbbot"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"github.com/kelseyhightower/envconfig"
	"regexp"
	"strings"
)

type MyConfigure struct {
	PAGEACCESSTOKEN string
	VERIFYTOKEN     string
	PORT 			int

}
var bot MyConfigure

type MyDatabase struct {
	NAME string
	USER string
	PASS string
}
var mydatabase MyDatabase

type TutorialUrl struct {
	URL string
}

var tutorialurl TutorialUrl

type Constant struct {
	MAXSAMPLE int
	MAXINNING int
}
var constant Constant

var db *sql.DB
type Record struct {}

func Init()  {
	err := envconfig.Process("BOT", &bot)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(bot.PAGEACCESSTOKEN,bot.VERIFYTOKEN,bot.PORT)

	err1 := envconfig.Process("DB", &mydatabase)
	if err1 != nil {
		log.Fatal(err.Error())
	}
	log.Println(mydatabase.NAME,mydatabase.USER,mydatabase.PASS)

	err2 := envconfig.Process("CONSTANT", &constant)
	if err2 != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(constant.MAXINNING,constant.MAXSAMPLE)

	err3 := envconfig.Process("URL", &tutorialurl)
	if err3 != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(tutorialurl.URL)
}

func ( record Record ) HandleMessage( bot *fbbot.Bot , msg *fbbot.Message ){
	if IsNewUser(db,msg.Sender.ID) {
		greeting := "Xin chào "+msg.Sender.FirstName()+", Chúng tôi đang thực hiện một dự án thu thập dữ liệu ghi âm giọng nói và rất vui khi nhận được sự hợp tác của bạn"
		m := fbbot.NewTextMessage(greeting)
		bot.Send(msg.Sender,m)

		tutorialmesseger :="Bạn hãy đọc và ghi âm lại các câu văn do chúng tôi gửi!"
		m1 := fbbot.NewTextMessage(tutorialmesseger)
		bot.Send(msg.Sender,m1)


		_,err := db.Query("INSERT INTO UserState(FbId,LastSample,Inning) VALUES(?,?,?)",msg.Sender.ID,0,1)
		if err!=nil{
			log.Println("error when insertNewUser")
		}
		btnms :=fbbot.NewButtonMessage()
		btnms.AddWebURLButton("Hướng dẫn",tutorialurl.URL)
		btnms.AddPostbackButton("Cập nhật thông tin","capnhatthongtin")
		btnms.Text="Bạn muốn xem hướng dẫn hay hoàn thành quá trình hoàn thành thông tin cho quá trình đăng ký"
		bot.Send(msg.Sender,btnms)

	}else if v,match := Isprovince(msg.Text); match && !IsAudioMessage(msg)  {
		_,err := db.Query("UPDATE UserState SET Province=? WHERE FbId= ? ",v,msg.Sender.ID)
		if err != nil {
			log.Println("error when execute insert Province!!")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin tỉnh, thành phố")
			bot.Send(msg.Sender,m)
		}
		if !checkProvince(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Tỉnh của bạn nhập chưa đúng định dạng\n Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(msg.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Bạn đã sẵn sàng cập nhật thông tin về tỉnh và thành phố quê hương bạn"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(msg.Sender,b)
		} else if !checkName(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin về tên  của bạn"
			b.AddPostbackButton("OK","okten")
			bot.Send(msg.Sender,b)
		}else if !checkPhoneNumber(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện th oạicủa bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập số điện thoại"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(msg.Sender,b)
		} else if !checkAge(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(msg.Sender,b)
		}else {
			_, err:=db.Query("UPDATE UserState SET State=? WHERE FbId= ? ",true,msg.Sender.ID)
			if err != nil {
				log.Println("error when update State for UserState")
			}
			m := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình cập nhật thông tin")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text= "Click để bắt đầu quá trình ghi âm"
			b.AddPostbackButton("Bắt đầu","batdaughiam")
			bot.Send(msg.Sender,b)
		}


	}else if v,match:=IsName(msg.Text);match && !IsAudioMessage(msg) {
		_,err := db.Query("UPDATE UserState SET Name=? WHERE FbId= ? ",v,msg.Sender.ID)
		if err != nil {
			log.Println("error when execute insert Name for user!!")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin về tên")
			bot.Send(msg.Sender,m)
		}
		if !checkName(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Tên bạn nhập chưa đúng định dạng \nChú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","okten")
			bot.Send(msg.Sender,b)
		}else if !checkProvince(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(msg.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text= "Bạn hãy cập nhật thông tin về tỉnh thành phố quê hương bạn"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(msg.Sender,b)
		} else  if !checkPhoneNumber(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện thoại của bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin số điện thoại"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(msg.Sender,b)
		} else if !checkAge(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng cập nhật thông tin về tuổi của bạn"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(msg.Sender,b)
		}else {
			_, err:=db.Query("UPDATE UserState SET State=? WHERE FbId= ? ",true,msg.Sender.ID)
			if err != nil {
				log.Println("error when update State for UserState")
			}
			m := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình cập nhật thông tin")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text= "Click để bắt đầu quá trình ghi âm"
			b.AddPostbackButton("Bắt đầu","batdaughiam")
			bot.Send(msg.Sender,b)
		}

	}else if v, match:=IsPhoneNumber(msg.Text);match && !IsAudioMessage(msg) {
		_,err := db.Query("UPDATE UserState SET NumberPhone=? WHERE FbId= ? ",v,msg.Sender.ID)
		if err != nil {
			log.Println("error when execute insert numberphone!!")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin số điện thoại")
			bot.Send(msg.Sender,m)
		}
		if !checkPhoneNumber(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Số điện thoại bạn vừa cập nhật chưa đúng định dạng! \n Chú ý: Số điện thoại của bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin về số điện thoại"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(msg.Sender,b)
		} else if !checkProvince(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(msg.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text = "Bạn hiểu và sẵn sàng cập nhật thông tin về tỉnh, thành phố của bạn?"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(msg.Sender,b)
		} else if !checkName(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","okten")
			bot.Send(msg.Sender,b)
		}else  if !checkAge(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng cập nhật thông tin tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(msg.Sender,b)
		}else {
			_, err:=db.Query("UPDATE UserState SET State=? WHERE FbId= ? ",true,msg.Sender.ID)
			if err != nil {
				log.Println("error when update State for UserState")
			}
			m := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình cập nhật thông tin")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text= "Click để bắt đầu quá trình ghi âm"
			b.AddPostbackButton("Bắt đầu","batdaughiam")
			bot.Send(msg.Sender,b)
		}

	}else if v,match:=IsAge(msg.Text);match && !IsAudioMessage(msg) {
		_,err := db.Query("UPDATE UserState SET Age=? WHERE FbId= ? ",v,msg.Sender.ID)
		if err != nil {
			log.Println("error when execute insert Age!!")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin về tuổi")
			bot.Send(msg.Sender,m)
		}
		if !checkAge(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("\n Tuổi bạn nhập không đúng định dạng\nChú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập lại thông tin tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(msg.Sender,b)
		}else if !checkProvince(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(msg.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Bạn đã rõ và sẵn sàng nhập thông tin tỉnh"
			b.AddPostbackButton("Click vào đây để tiếp tục","oktinh")
			bot.Send(msg.Sender,b)
		} else if !checkName(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","okten")
			bot.Send(msg.Sender,b)
		}else if !checkPhoneNumber(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện th oạicủa bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(msg.Sender,b)
		}else {
			_, err:=db.Query("UPDATE UserState SET State=? WHERE FbId= ? ",true,msg.Sender.ID)
			if err != nil {
				log.Println("error when update State for UserState")
			}
			m := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình cập nhật thông tin")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text= "Click để bắt đầu quá trình ghi âm"
			b.AddPostbackButton("Bắt đầu","batdaughiam")
			bot.Send(msg.Sender,b)
		}
	}else if !AvailableUser(db,msg.Sender.ID) {
		if !checkProvince(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(msg.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Mời bạn cập nhật thông tin tỉnh, thành phố"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(msg.Sender,b)
		} else if !checkName(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","okten")
			bot.Send(msg.Sender,b)
		}else if !checkPhoneNumber(db,msg.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện thoại của bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(msg.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập số điện thoại"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(msg.Sender,b)
		} else if !checkAge(db,msg.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(msg.Sender,b)
		} else {
			_, err:=db.Query("UPDATE UserState SET State=? WHERE FbId= ? ",true,msg.Sender.ID)
			if err != nil {
				log.Println("error when update State for UserState")
			}
			m := fbbot.NewTextMessage("Bạn đã hoàn thành quá trình cập nhật thông tin")
			bot.Send(msg.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text= "Click để bắt đầu quá trình ghi âm"
			b.AddPostbackButton("Bắt đầu","batdaughiam")
			bot.Send(msg.Sender,b)
		}
	} else if IsAudioMessage(msg) {
		log.Println("is audio!!!!")
		st := GetCurrentState(db, msg.Sender.ID)
		id := st + 1
		inning := GetCurrentInning(db, msg.Sender.ID)
		if !isExist(db, msg.Sender.ID, id, inning) {
			_, err := db.Query("INSERT INTO Outputs(FbId, SampleId, State, Inning, UrlRecord) Value(?,?,?,?,?)", msg.Sender.ID, id, false, inning, msg.Audios[0].URL)
			if err != nil {
				log.Println("error when Insert to outputs")
			}
		} else {
			state := GetCurrentState(db, msg.Sender.ID)
			sampleid := state + 1
			stmtInsAudio, err := db.Prepare("UPDATE Outputs SET UrlRecord=? WHERE FbId= ? AND SampleId=?")
			if err != nil {
				log.Println("error when create stminsertAudio")
			}
			stmtInsAudio.Query(msg.Audios[0].URL, msg.Sender.ID, sampleid)
		}
		btnms := fbbot.NewButtonMessage()
		btnms.AddPostbackButton("Ghi âm lại", "ghiamlai")
		btnms.AddPostbackButton("Ghi âm câu tiếp theo", "cautieptheo")
		btnms.Text = "Bạn muốn ghi âm lại hay ghi âm câu tiếp theo"
		bot.Send(msg.Sender, btnms)
	}else {
		state := GetCurrentState(db,msg.Sender.ID)
		smlid := state+1
		inning := GetCurrentInning(db,msg.Sender.ID)
		if inning<= constant.MAXINNING {
			if smlid <= constant.MAXSAMPLE{
				t := fbbot.NewTextMessage("Mời bạn thu âm câu sau")
				bot.Send(msg.Sender,t)
				sample := GetSample(db,smlid)
				m := fbbot.NewTextMessage(sample)
				bot.Send(msg.Sender,m)
			}
		}else {
			sample:= "Bạn đã hoàn thành quá trình ghi âm, xin chân thành cảm ơn"
			m:= fbbot.NewTextMessage(sample)
			bot.Send(msg.Sender,m)
		}
	}
}

func (r Record) HandlePostback(bot *fbbot.Bot, pbk *fbbot.Postback)  {
	switch pbk.Payload {
	case "capnhatthongtin":
		provincechooser := fbbot.NewButtonMessage()
		provincechooser.AddPostbackButton("Miền Bắc","mienbac")
		provincechooser.AddPostbackButton("Miền Trung","mientrung")
		provincechooser.AddPostbackButton("Miền Nam","miennam")
		provincechooser.Text = "Bạn thuộc vùng miền nào của Việt Nam"
		bot.Send(pbk.Sender,provincechooser)
	case "batdaughiam":
			state := GetCurrentState(db,pbk.Sender.ID)
			smlid := state+1
			inning := GetCurrentInning(db,pbk.Sender.ID)
			if inning<= constant.MAXINNING {
				if smlid <= constant.MAXSAMPLE{
					t := fbbot.NewTextMessage("Mời bạn thu âm câu sau")
					bot.Send(pbk.Sender,t)
					sample := GetSample(db,smlid)
					m := fbbot.NewTextMessage(sample)
					bot.Send(pbk.Sender,m)
				}
			}else {
				sample:= "Bạn đã hoàn thành quá trình ghi âm, xin chân thành cảm ơn"
				m:= fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)
			}

	case "mienbac":
		_ , err :=db.Query("UPDATE UserState SET Area =? WHERE FbId= ? ","Miền Bắc",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updatevungmien user")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin vùng miền")
			bot.Send(pbk.Sender,m)
		}

		if !checkProvince(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(pbk.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Bạn đã sẵn sàng cập nhật thông tin về tỉnh thành phố của bạn?"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(pbk.Sender,b)
		} else if !checkName(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin về tên của bạn"
			b.AddPostbackButton("OK","okten")
			bot.Send(pbk.Sender,b)
		}else if !checkPhoneNumber(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện th oạicủa bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(pbk.Sender,b)
		} else if !checkAge(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(pbk.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(pbk.Sender,b)
		}
	case "mientrung":
		_ , err :=db.Query("UPDATE UserState SET Area =? WHERE FbId= ? ","Miền Trung",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updatevungmien user")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin vùng miền")
			bot.Send(pbk.Sender,m)
		}

		if !checkProvince(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(pbk.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Bạn đã sẵn sàng cập nhật thông tin về tỉnh thành phố của bạn?"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(pbk.Sender,b)
		} else if !checkName(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng: TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin về tên của bạn"
			b.AddPostbackButton("OK","okten")
			bot.Send(pbk.Sender,b)
		}else if !checkPhoneNumber(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện th oạicủa bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(pbk.Sender,b)
		} else if !checkAge(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(pbk.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(pbk.Sender,b)
		}
	case "miennam":
		_ , err :=db.Query("UPDATE UserState SET Area =? WHERE FbId= ? ","Miền Nam",pbk.Sender.ID )
		if err != nil {
			log.Println("error when execute updatevungmien user")
		}else {
			m := fbbot.NewTextMessage("Bạn đã cập nhật xong thông tin vùng miền")
			bot.Send(pbk.Sender,m)
		}
		if !checkProvince(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tên tỉnh viết liền không dấu (VD đúng :ThaiBinh)(VD sai:tỉnh Thái Bình")
			bot.Send(pbk.Sender,m)

			b :=fbbot.NewButtonMessage()
			b.Text="Bạn đã sẵn sàng cập nhật thông tin về tỉnh thành phố của bạn?"
			b.AddPostbackButton("OK","oktinh")
			bot.Send(pbk.Sender,b)
		} else if !checkName(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Tên của bạn phải theo định dạng (VD đúng:TênTôiLà: Nguyễn Việt Tiến)\n(VD sai: tên của tôi là Nguyễn Việt Tiến")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng cập nhật thông tin về tên của bạn"
			b.AddPostbackButton("OK","okten")
			bot.Send(pbk.Sender,b)
		}else if !checkPhoneNumber(db,pbk.Sender.ID) {
			m :=fbbot.NewTextMessage("Chú ý: Số điện th oạicủa bạn phải theo định dạng (VD đúng: 0974793322) (VD sai: số điện thoại của tôi là +84974793322")
			bot.Send(pbk.Sender,m)
			b :=fbbot.NewButtonMessage()
			b.Text ="Bạn đã hiểu và sẵn sàng nhập tên"
			b.AddPostbackButton("OK","oksodienthoai")
			bot.Send(pbk.Sender,b)
		} else if !checkAge(db,pbk.Sender.ID){
			m := fbbot.NewTextMessage("Chú ý: Tuổi của bạn phải đúng định dạng (VD đúng: 21) (VD sai: 21 tuoi)")
			bot.Send(pbk.Sender,m)
			b := fbbot.NewButtonMessage()
			b.Text = "Bạn đã hiểu và sẵn sàng nhập tuổi"
			b.AddPostbackButton("OK","oktuoi")
			bot.Send(pbk.Sender,b)
		}
	case "cautieptheo":
		state := GetCurrentState(db,pbk.Sender.ID)
		smlid :=state+1
		ig    := GetCurrentInning(db,pbk.Sender.ID)
		if ig<= constant.MAXINNING {
			if smlid <constant.MAXSAMPLE {
				_, err := db.Query("UPDATE Outputs SET State = ? WHERE FbId=? AND SampleId=? AND Inning=?",true,pbk.Sender.ID,smlid,ig)
				if err != nil {
					log.Println("error when update state of outputs")
				}

				_, err1 := db.Query("UPDATE UserState SET LastSample=? WHERE FbId=? ",smlid,pbk.Sender.ID)
				if err1 != nil {
					log.Println("error when update state of user ")
				}
				t := fbbot.NewTextMessage("Mời bạn thu âm câu sau")
				bot.Send(pbk.Sender,t)
				sample := GetSample(db,smlid+1)
				m := fbbot.NewTextMessage(sample)
				bot.Send(pbk.Sender,m)

			}else if smlid == constant.MAXSAMPLE {
				_, err := db.Query("UPDATE Outputs SET State=? WHERE FbId=? AND SampleId=? AND Inning=?",true,pbk.Sender.ID,smlid,ig)
				if err != nil {
					log.Println("error when update state of outputs")
				}
				_, err1 := db.Query("UPDATE UserState SET LastSample=?, Inning =? where FbId=? ",0,ig+1, pbk.Sender.ID)
				if err1 != nil {
					log.Println("error when update state of user ")
				}
				ann1 := "Bạn đã hoàn thành lượt:"
				ann2 := strconv.Itoa(ig)
				announce := ann1+ann2
				m := fbbot.NewTextMessage(announce)
				bot.Send(pbk.Sender,m)
				ig:=GetCurrentInning(db,pbk.Sender.ID)
				if ig>constant.MAXINNING{
					sample := "Bạn đã hoàn thành quá trình thu âm, cảm ơn bạn"
					m := fbbot.NewTextMessage(sample)
					bot.Send(pbk.Sender,m)
				} else {
					t := fbbot.NewTextMessage("Mời bạn thu âm câu sau ")
					bot.Send(pbk.Sender,t)
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
		if ig<=constant.MAXINNING {
			t := fbbot.NewTextMessage("Mời bạn thu âm câu "+string(sampleid) + " :")
			bot.Send(pbk.Sender,t)
			sample :=GetSample(db,sampleid)
			m:= fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}else {
			sample := "Bạn đã hoàn thành xong bài ghi âm"
			m:= fbbot.NewTextMessage(sample)
			bot.Send(pbk.Sender,m)
		}
	case "oktinh":
		t :="Mời bạn nhập tỉnh, thành phố quê hương của bạn!"
		m :=fbbot.NewTextMessage(t)
		bot.Send(pbk.Sender,m)
	case "okten":
		t :="Mời bạn nhập tên !"
		m :=fbbot.NewTextMessage(t)
		bot.Send(pbk.Sender,m)
	case "oksodienthoai":
		t :="Mời bạn nhập thông tin số điện thoại!"
		m :=fbbot.NewTextMessage(t)
		bot.Send(pbk.Sender,m)
	case "oktuoi":
		t :="Mời bạn nhập tuổi!"
		m :=fbbot.NewTextMessage(t)
		bot.Send(pbk.Sender,m)
	default:
		log.Println("no case in switch")
	}
}
func main() {
	Init()

	//processing database
	var err error
	db, err = sql.Open("mysql", mydatabase.USER+":"+mydatabase.PASS+"@/"+mydatabase.NAME )//"user:password@/dbname"
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
	bot := fbbot.New(bot.PORT,bot.VERIFYTOKEN,bot.PAGEACCESSTOKEN)
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
//todo
func IsPhoneNumber(pn string) (string,bool)  {
	r,err := regexp.Compile("0[0-9]{9,10}$")
	if err!= nil {
		log.Println("err compile phonenumber regex")
	}
	if r.FindString(pn)!="" {
		return r.FindString(pn),true
	}
	return r.FindString(pn),false
}
func IsAge(age string) (string, bool)  {
	if(len(age))>2 {
		return "",false
	}else {
		r, err := regexp.Compile("([0-9])[0-9]{0,1}$")
		if err!= nil {
			log.Println("err compile Age regex")
		}
		s := r.FindString(age)
		fmt.Println("IsAgefunc: age= ",s)
		if s=="99"{
			return "",false
		}else {
			return s,true
		}
	}
}
func Isprovince(pe string) (string,bool) {
	a := strings.ToLower(pe)
	b := strings.Split(a," ")
	v := strings.Join(b,"")
	r, err:= regexp.Compile("(tp|thanhpho|)(hanoi|hochiminh|angiang|bariavungtau|baclieu|bacgiang|backan|bacninh|bentre|binhduong|binhdinh|binhphuoc|binhthuan|binhthuan|caobang|cantho|haiphong|danang|gialai|hoabinh|hagiang|hanam|hatinh|hungyen|haiduong|haugiang|dienbien|daklak|daknong|lamdong|dongnai|dongthap|khanhhoa|kiengiang|kontum|laichau|longan|laocai|langson|namdinh|nghean|ninhbinh|ninhthuan|phutho|phuyen|quangbinh|quangnam|quangngai|quangninh|quangtri|soctrang|sonla|thanhhoa|thaibinh|thainguyen|thuathienhue|thuathien-hue|tiengiang|travinh|tuyenquang|tayninh|vinhlong|vinhphuc|yenbai)")
	if err != nil {
		log.Println("do not compile pattern")
	}
	if r.MatchString(v) {
		return v,true
	}else {
		return "",false
	}
}
func IsName(name string) (string, bool) {
	a :=strings.ToLower(name)
	b :=strings.Split(a," ")
	t :=strings.Join(b,"")
	r,err := regexp.Compile("(têntôilà):(.{1,})$")
	if err!=nil {
		fmt.Println("error when compile IsName")
	}
	if r.MatchString(t) {
		m :=strings.Replace(t,"têntôilà:","",-1)
		return m,true
	}else {
		return "", false
	}
}
func AvailableUser(db *sql.DB,FbId string) bool {
	var State bool
	row,err := db.Query("SELECT State from UserState WHERE FbId=?",FbId)
	if err!=nil {
		log.Println("err when query to get state of user")
	}
	for row.Next() {
		row.Scan(&State)
	}
	return State
}
func checkName(db *sql.DB, FbId string) bool {
	var name string
	row, err := db.Query("SELECT Name from UserState WHERE FbId = ?",FbId)
	if err!=nil {
			log.Println("err when checkName")
	}
	row.Next()
	row.Scan(&name)
	fmt.Println("tên: ",name)
	if name=="empty" {
		return false
	}else {
		return true
	}
}
func checkPhoneNumber(db *sql.DB,FbId string) bool {
	var np string
	row, err := db.Query("SELECT NumberPhone FROM UserState WHERE FbId=?",FbId)
	if err!=nil {
		log.Println("err when query check NumberPhone")
	}
	row.Next()
	row.Scan(&np)
	fmt.Println("số điện thoại",np)
	if np=="empty" {
		return false
	} else {
		return true
	}

}
func checkProvince(db *sql.DB,FbId string) bool {
	var province string
	row , err := db.Query("SELECT Province FROM UserState WHERE FbId = ?", FbId)
	if err!=nil {
		log.Println("error when query checkProvince")
	}
	row.Next()
	row.Scan(&province)
	fmt.Println("tên tỉnh:",province)
	if province=="empty" {
		return false
	} else {
		return true
	}
}

func checkAge(db *sql.DB,FbId string) bool {
	var age string
	row , err := db.Query("SELECT Age FROM UserState WHERE FbId = ?", FbId)
	if err!=nil {
		log.Println("error when query checkAge")
	}
	row.Next()
	row.Scan(&age)
	fmt.Println("tuổi: ",age)
	if age=="99" {
		return false
	} else {
		return true
	}
}
