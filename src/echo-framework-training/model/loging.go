package model

type LoginData struct {
	Email    string `json:"mail"`
	PassWord string `json:"pass"`
}

type UserInfo struct {
	UserID     int    `json:"user_id"`
	UserCode   string `json:"user_code"`
	PassWord   string `json:"pass"`
	Mail       string `json:"mail"`
	TalkPoint  int    `json:"talk_pt"`
	ChekiPoint int    `json:"cheki_pt"`
	SignPoint  int    `json:"sign_pt"`
}

func UserLogIn(info *LoginData) (interface{}, bool) {
	// テーブル検索してあればページ遷移
	user := UserInfo{}
	db := ConnectDB()
	println(info.Email)
	row := db.QueryRow("select user_id, user_code, password, mail, talk_pt, cheki_pt, sign_pt from db_server.usr_table where mail=\"" + info.Email + "\" && password=\"" + info.PassWord + "\"")
	text := "select user_id, user_code, usr_name, password, mail from db_server.usr_table where mail=\"" + info.Email + "\" && password=\"" + info.PassWord + "\""
	println(text)
	if err := row.Scan(&user.UserID, &user.UserCode, &user.PassWord, &user.Mail, &user.TalkPoint, &user.ChekiPoint, &user.SignPoint); err != nil {
		// ここでエラー出てる
		println(err.Error())
	}
	println(user.Mail)
	result := user.Mail != "" || user.UserID != 0
	return user, result
}
