package model

import "strconv"

type ReserveInfo struct {
	RecordID     string `json:"id"`
	UserRecordID string `json:"record_id"`
	UserID       string `json:"user_id"`
	UserCode     string `json:"user_code"`
	Member       string `json:"member"`
	Venue        string `json:"venue"`
	Type         string `json:"type"`
	Department   string `json:"department"`
}

type RecordNumber struct {
	UserRecordID int `json:"user_record_id"`
}

// 予約する際に呼ばれるAPI
func ReserveTicket(info *ReserveInfo) string {
	// 予約処理
	db := ConnectDB()
	master := GetMasterData()

	var member_id string
	for i := 0; i < len(master.MemberMaster); i++ {
		if info.Member == master.MemberMaster[i].Name {
			member_id = master.MemberMaster[i].MemberId
		}
	}
	var venue_id string
	for i := 0; i < len(master.VenueMaster); i++ {
		if info.Venue == master.VenueMaster[i].Venue_STR {
			venue_id = master.VenueMaster[i].VenueId
		}
	}
	var type_id string
	for i := 0; i < len(master.TypeMaster); i++ {
		if info.Type == master.TypeMaster[i].Type_STR {
			type_id = master.TypeMaster[i].TypeId
		}
	}
	var department_id string
	for i := 0; i < len(master.DepartmentMaster); i++ {
		if info.Department == master.DepartmentMaster[i].Department_STR {
			department_id = master.DepartmentMaster[i].DepartmentId
		}
	}
	// TODO ユーザーレコードのIDを取得
	User_ID := "100000"
	User_Code := "'v36zc93cuicx'"
	rowQuery := "select * from reserve_table where user_id=" + User_ID + " && member_id=" + member_id + " && venue_id=" + venue_id + " && type_id=" + type_id + " && department_id=" + department_id
	println(rowQuery)
	row := db.QueryRow(rowQuery)
	exist := ReserveInfo{}
	row.Scan(&exist.RecordID, &exist.UserRecordID, &exist.UserID, &exist.UserCode, &exist.Member, &exist.Venue, &exist.Type, &exist.Department)
	if exist.UserRecordID != "" && exist.UserID != "" && exist.UserCode != "" {
		// すでに存在しているレコードなので予約不要
		return "already"
	}
	println("record id=" + exist.UserCode + " UserID=" + exist.UserID + " UserCode=" + exist.UserCode)
	nmb := db.QueryRow("select user_record_id from reserve_table where user_id=? order by user_record_id DESC", "100000")
	number := RecordNumber{}
	nmb.Scan(&number.UserRecordID)
	record_id := 1
	println(number.UserRecordID)
	if number.UserRecordID != 0 {
		record_id = number.UserRecordID + 1
	}
	query := "insert into db_server.reserve_table (user_record_id, user_id, user_code, member_id, venue_id, type_id, department_id) values (" + strconv.Itoa(record_id) + "," + User_ID + "," + User_Code + "," + member_id + "," + venue_id + "," + type_id + "," + department_id + ")"
	println(query)
	ins, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	ins.Exec()
	println("実行しました")
	defer db.Close()
	println("リターンを返します。")
	return "done"
}
