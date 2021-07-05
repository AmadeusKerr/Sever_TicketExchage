package model

type UserRecordId struct {
	RecordID     int `json:"record_id"`
	UserRecordNO int `json:"user_record_no"`
}

// 予約一覧から消す処理
func RemoveReservedRocord(reserveId, recordNo string, userId int) bool {
	db := ConnectDB()
	delForm, err := db.Prepare("DELETE FROM reserve_table WHERE reserve_id=? && user_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(reserveId, userId)

	rows, err := db.Query("select reserve_id, user_record_id from reserve_table where user_id=?", userId)
	if err != nil {
		panic(err.Error())
	}

	count := 1
	for rows.Next() {
		recordId := UserRecordId{}
		if err := rows.Scan(&recordId.RecordID, &recordId.UserRecordNO); err != nil {
			panic(err.Error())
		}
		number, err := db.Prepare("update reserve_table set user_record_id=? where user_id=? && reserve_id=?")
		if err != nil {
			panic(err.Error())
		}
		number.Exec(count, userId, recordId.RecordID)
		count++
	}
	defer db.Close()
	return true
}
