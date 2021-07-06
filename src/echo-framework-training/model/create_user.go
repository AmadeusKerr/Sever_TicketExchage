package model

import (
	"math/rand"
	"strconv"
)

type CreateUserInfo struct {
	Mail          string `json:"mail"`
	PassWord      string `json:"pass"`
	UserNameFirst string `json:"user_name_first"`
	UserNameLast  string `json:"user_name_last"`
	PostCode      string `json:"post_code"`
	Address       string `json:"address"`
	Phone         string `json:"phone_number"`
}

type UsesrIdHolder struct {
	UserId string
}

func CreateUser(info *CreateUserInfo) interface{} {
	// UserIdの検索 => 0件 => 作成(ラストインデックス+1)
	db := ConnectDB()
	u := new(UsesrIdHolder)
	check := db.QueryRow("select user_id from db_server.usr_table where mail=" + ConvertSQLString(info.Mail))
	check.Scan(&u.UserId)
	println(u.UserId)
	if u.UserId != "" && u.UserId != "0" {
		return map[string]string{"result": "already"}
	}
	user := new(UsesrIdHolder)
	code := GenerateUserCode()
	println(code)
	row := db.QueryRow("SELECT user_id FROM db_server.usr_table ORDER BY user_id DESC LIMIT 1")
	row.Scan(&user.UserId)
	userIdInt, err := strconv.Atoi(user.UserId)
	if err != nil {
		panic(err.Error())
	}
	id := strconv.Itoa((userIdInt + 1))
	user.UserId = id

	query := "insert into db_server.usr_table (user_code, first_name, last_name, password, mail, postal_code, address, phone_number, talk_pt, cheki_pt, sign_pt, create_at, update_at) values (?,?,?,?,?,?,?,?,?,?)"
	ins, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	_, err = ins.Exec(
		ConvertSQLString(code),
		NewNullString(info.UserNameFirst),
		NewNullString(info.UserNameLast),
		ConvertSQLString(info.PassWord),
		ConvertSQLString(info.Mail),
		NewNullString(info.PostCode),
		NewNullString(info.Address),
		NewNullString(info.Phone),
		0, 0, 0,
		ConvertSQLString(GetDate()),
		ConvertSQLString(GetDate()))
	if err != nil {
		println(err.Error())
	}
	defer db.Close()
	userMap := make(map[string]interface{})
	userMap["user_id"] = user.UserId
	userMap["user_code"] = code
	userMap["result"] = "done"

	// ローカルに保持
	return userMap // mapにして返した方がいいかも
}

func GenerateUserCode() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 11)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
